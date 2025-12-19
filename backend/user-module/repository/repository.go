package repository

import (
	"context"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/user-module/model"

	"gorm.io/gorm"
)

// IUserRepository defines the interface for user repository
type IUserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *model.User) error
	GetByID(ctx context.Context, tx *gorm.DB, id int) (*model.User, error)
	Update(ctx context.Context, tx *gorm.DB, user *model.User) error
	Delete(ctx context.Context, tx *gorm.DB, id int) error
}

// Repository implements IUserRepository
type Repository struct {
	DB     *gorm.DB
	Logger LoggerInterface
}

func (r Repository) Create(ctx context.Context, tx *gorm.DB, user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) GetByID(ctx context.Context, tx *gorm.DB, id int) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Update(ctx context.Context, tx *gorm.DB, user *model.User) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository) Delete(ctx context.Context, tx *gorm.DB, id int) error {
	//TODO implement me
	panic("implement me")
}

type LoggerInterface interface {
	Info(args ...interface{})
	Error(args ...interface{})
}
