package main

import (
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
	"net/http"
	"time"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/core"
	httpModule "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/server/http"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Create a new router
	router := mux.NewRouter()

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

	// Initialize core
	userCore, err := core.NewCore(db, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize core: %v", err)
	}

	// Create HTTP extension
	httpExt := httpModule.NewUserHttpExtension(router, userCore)

	// Register all routes
	httpModule.RegisterRoutes(httpExt)

	// Start the server
	logger.Info("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
