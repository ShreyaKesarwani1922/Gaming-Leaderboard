package core

import (
	"fmt"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
	userDataMapper "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/datamapper"
	userRepository "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/repository"
	"gorm.io/gorm"
)

type ICore interface {
}

type Core struct {
	DB         *gorm.DB
	Repository *userRepository.IUserRepository
	DataMapper *userDataMapper.DataMapper
	Logger     *providers.LoggerInterface
}

func NewCore(db *gorm.DB, logger providers.LoggerInterface) (*Core, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection cannot be nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	repo := userRepository.NewRepository(db, logger)

	// Initialize data mapper
	//dataMapper := userDataMapper.NewDataMapper()

	return &Core{
		DB:         db,
		Repository: &repo,
		//DataMapper: dataMapper,
		Logger: &logger,
	}, nil
}
