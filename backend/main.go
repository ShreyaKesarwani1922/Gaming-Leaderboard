package main

import (
	"net/http"
	"time"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"

	dataMigrationCore "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/core"
	dataMigrationRepo "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/repository"
	dataMigrationHttpModule "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/server/http"
	leaderBoardCore "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/core"
	leaderBoardRepo "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/repository"
	leaderBoardHttp "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/server/http"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/core"
	httpModule "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/server/http"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// ------------------------------------------------------------------
	// Logger
	// ------------------------------------------------------------------
	logger := providers.NewConsoleLogger()
	logger.Info("Starting Gaming-Leaderboard Backend Service")

	// ------------------------------------------------------------------
	// Router
	// ------------------------------------------------------------------
	router := mux.NewRouter()

	// Panic recovery middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Errorf("Panic recovered | path=%s error=%v", r.URL.Path, rec)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})

	// Request logging + New Relic transaction naming
	router.Use(func(next http.Handler) http.Handler {
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
	})

	// ------------------------------------------------------------------
	// Database
	// ------------------------------------------------------------------
	logger.Info("Connecting to PostgreSQL database")

	dsn := "host=localhost user=postgres password=password dbname=gaming_dashboard port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("Failed to get sql.DB instance: %v", err)
	}

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
		newrelic.ConfigLicense("6366ed4fc1a70027f322f15ec705de53FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		logger.Fatalf("New Relic initialization failed: %v", err)
	}

	logger.Info("New Relic initialized successfully")

	// ------------------------------------------------------------------
	// User Module
	// ------------------------------------------------------------------
	logger.Info("Initializing User module")

	userCore, err := core.NewCore(db, logger)
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

	leaderboardRepo := leaderBoardRepo.NewLeaderBoardRepository(db, logger)
	leaderboardCore := leaderBoardCore.NewLeaderboardCore(leaderboardRepo, logger)
	leaderboardHandler := leaderBoardHttp.NewLeaderboardHandler(leaderboardCore, logger, nrApp)
	leaderboardHandler.RegisterRoutes(router)

	logger.Info("Leaderboard routes registered")

	// ------------------------------------------------------------------
	// Start Server
	// ------------------------------------------------------------------
	port := ":3000"
	logger.Infof("Server listening on %s", port)

	if err := http.ListenAndServe(port, router); err != nil {
		logger.Fatalf("Server startup failed: %v", err)
	}
}
