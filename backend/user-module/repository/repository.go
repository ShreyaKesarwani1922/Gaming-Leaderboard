package repository

import (
	"fmt"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/model"
	"gorm.io/gorm"
)

// IUserRepository defines the interface for user repository
type IUserRepository interface{}

// Repository implements IUserRepository
type Repository struct {
	DB     *gorm.DB
	Logger *providers.LoggerInterface
}

// NewRepository creates a new instance of the user repository
func NewRepository(db *gorm.DB, logger providers.LoggerInterface) IUserRepository {
	if db == nil {
		panic("database connection cannot be nil")
	}

	if logger == nil {
		panic("logger cannot be nil")
	}

	// Auto-migrate the User model
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to auto-migrate User model: %v", err))
		panic(fmt.Sprintf("failed to auto-migrate User model: %v", err))
	}

	return &Repository{
		DB:     db,
		Logger: &logger,
	}
}
