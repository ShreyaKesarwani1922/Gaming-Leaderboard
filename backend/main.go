package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	dataMigrationCore "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/core"
	dataMigrationRepo "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/repository"
	dataMigrationHttpModule "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/server/http"
	leaderBoardCore "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/core"
	leaderBoardRepo "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/repository"
	leaderBoardHttp "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/server/http"
	userCore "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/core"
	httpModule "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/server/http"
)

func main() {
	// ------------------------------------------------------------------
	// Logger
	// ------------------------------------------------------------------
	logger := providers.NewConsoleLogger()
	logger.Info("Starting Gaming-Leaderboard Backend Service")

	// ------------------------------------------------------------------
	// Configuration
	// ------------------------------------------------------------------
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0 // default DB

	// ------------------------------------------------------------------
	// Redis Client
	// ------------------------------------------------------------------
	logger.Info("Initializing Redis client", "address", redisAddr)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
		PoolSize: 20, // Adjust based on your needs
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	logger.Info("Redis connected successfully")

	// ------------------------------------------------------------------
	// Router
	// ------------------------------------------------------------------
	router := mux.NewRouter()

	// Panic recovery middleware
	router.Use(panicRecovery(logger))
	router.Use(requestLogger(logger))

	// ------------------------------------------------------------------
	// Database
	// ------------------------------------------------------------------
	logger.Info("Connecting to PostgreSQL database")

	dsn := getEnv("DATABASE_URL", "host=localhost user=postgres password=password dbname=gaming_dashboard port=5432 sslmode=disable TimeZone=UTC")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("Failed to get sql.DB instance: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("PostgreSQL connected and pool configured")

	// ------------------------------------------------------------------
	// New Relic
	// ------------------------------------------------------------------
	logger.Info("Initializing New Relic")

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Gaming-Leaderboard-Backend"),
		newrelic.ConfigLicense("1d04f778316b658766e6028835a8c0e6FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		logger.Warnf("New Relic initialization failed (continuing without it): %v", err)
	} else {
		logger.Info("New Relic initialized successfully")
	}

	// ------------------------------------------------------------------
	// User Module
	// ------------------------------------------------------------------
	logger.Info("Initializing User module")

	userCore, err := userCore.NewCore(db, logger)
	if err != nil {
		logger.Fatalf("User core initialization failed: %v", err)
	}

	userHttpExt := httpModule.NewUserHttpExtension(router, userCore)
	httpModule.RegisterRoutes(userHttpExt)

	logger.Info("User routes registered")

	// ------------------------------------------------------------------
	// Data Migration Module
	// ------------------------------------------------------------------
	logger.Info("Initializing Data Migration module")

	dmRepo := dataMigrationRepo.NewMigrationRepository(db)
	dmCore := dataMigrationCore.NewMigrationCore(dmRepo, db)
	dmHandler := dataMigrationHttpModule.NewMigrationHandler(dmCore, nrApp, logger)
	dmHandler.RegisterRoutes(router)

	logger.Info("Data Migration routes registered")

	// ------------------------------------------------------------------
	// Leader Board Module
	// ------------------------------------------------------------------
	logger.Info("Initializing Leader Board module")

	leaderboardRepo := leaderBoardRepo.NewLeaderBoardRepository(db, redisClient, logger)
	leaderboardCore := leaderBoardCore.NewLeaderboardCore(leaderboardRepo, logger)
	leaderboardHandler := leaderBoardHttp.NewLeaderboardHandler(leaderboardCore, logger, nrApp)
	leaderboardHandler.RegisterRoutes(router)

	logger.Info("Leaderboard routes registered")

	// ------------------------------------------------------------------
	// Health Check
	// ------------------------------------------------------------------
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}).Methods("GET")

	// ------------------------------------------------------------------
	// Start Server
	// ------------------------------------------------------------------
	port := getEnv("PORT", "8000")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Infof("Server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Server startup failed: %v", err)
	}
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// panicRecovery middleware handles panics and logs them
func panicRecovery(logger *providers.ConsoleLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Errorf("Panic recovered | path=%s error=%v", r.URL.Path, rec)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// requestLogger middleware logs incoming requests
func requestLogger(logger *providers.ConsoleLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			if txn := newrelic.FromContext(r.Context()); txn != nil {
				txn.SetName(r.Method + " " + r.URL.Path)
			}

			logger.Infof(
				"Incoming request | method=%s path=%s remote=%s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
			)

			next.ServeHTTP(w, r)

			logger.Infof(
				"Request completed | method=%s path=%s duration=%s",
				r.Method,
				r.URL.Path,
				time.Since(start),
			)
		})
	}
}
