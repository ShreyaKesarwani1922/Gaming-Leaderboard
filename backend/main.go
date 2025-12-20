package main

import (
	"net/http"
	"time"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"

	dataMigrationCore "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/core"
	dataMigrationRepo "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/repository"
	dataMigrationHttpModule "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/server/http"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/core"
	httpModule "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/server/http"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Create a new router
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if txn := newrelic.FromContext(r.Context()); txn != nil {
				txn.SetName(r.Method + " " + r.URL.Path)
			}
			next.ServeHTTP(w, r)
		})
	})

	// Initialize logger first
	logger := providers.NewConsoleLogger()

	// Initialize PostgreSQL database connection
	dsn := "host=localhost user=postgres password=password dbname=gaming_dashboard port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the database connection
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("Failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Successfully connected to PostgreSQL database")

	// Adding New Relic
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName("Gaming-Leaderboard-Backend"),
		newrelic.ConfigLicense("6366ed4fc1a70027f322f15ec705de53FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		logger.Fatalf("Failed to initialize New Relic: %v", err)
	}

	// Initialize userCore etc
	userCore, err := core.NewCore(db, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize core: %v", err)
	}
	httpExt := httpModule.NewUserHttpExtension(router, userCore)
	httpModule.RegisterRoutes(httpExt)

	// Initialize data migration
	dmRepo := dataMigrationRepo.NewMigrationRepository(db)
	dmCore := dataMigrationCore.NewMigrationCore(dmRepo, db)
	dmHandler := dataMigrationHttpModule.NewMigrationHandler(dmCore, nrApp)
	dmHandler.RegisterRoutes(router) // Call the method on the handler instance

	// Create a new router that will be wrapped by New Relic
	nrRouter := http.NewServeMux()
	nrRouter.HandleFunc(newrelic.WrapHandleFunc(nrApp, "/", func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	}))

	logger.Info("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
