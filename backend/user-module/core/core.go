package core

import (
	userDataMapper "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/datamapper"
	userRepository "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/repository"
	"gorm.io/gorm"
)

type ICore interface {
}

type Core struct {
	DB         *gorm.DB
	Repository *userRepository.Repository
	DataMapper *userDataMapper.DataMapper
	Logger     LoggerInterface
}
