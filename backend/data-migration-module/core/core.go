package core

import (
	"context"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/repository"
	"gorm.io/gorm"
)

type IMigrationCore interface {
	PopulateSampleData() error
}

type MigrationCore struct {
	repository *repository.MigrationRepository
	db         *gorm.DB
}

func NewMigrationCore(repo *repository.MigrationRepository, db *gorm.DB) *MigrationCore {
	return &MigrationCore{
		repository: repo,
		db:         db,
	}
}

func (c *MigrationCore) PopulateSampleData(ctx context.Context, userLimit, sessionLimit int) error {
	tx := c.db.Begin()
	if userLimit > 0 {
		if err := c.repository.BulkInsertUsers(tx, userLimit); err != nil {
			return err
		}
	}

	if sessionLimit > 0 {
		if err := c.repository.BulkInsertGameSessions(tx, sessionLimit); err != nil {
			return err
		}
	}

	if err := c.repository.UpdateLeaderboard(tx); err != nil {
		return err
	}

	return tx.Commit().Error
}
