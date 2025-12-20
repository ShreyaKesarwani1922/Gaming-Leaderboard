// repository/repository.go
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/constants"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
	"gorm.io/gorm"
)

type LeaderboardEntry struct {
	UserID     int64 `gorm:"column:user_id"`
	TotalScore int64 `gorm:"column:total_score"`
}

type ILeaderboardRepository interface {
	SubmitScore(ctx context.Context, userID, score int64, gameMode string) (time.Time, error)
	GetTopPlayers(ctx context.Context, limit int) ([]LeaderboardEntry, error)
}

type LeaderboardRepository struct {
	db     *gorm.DB
	logger *providers.ConsoleLogger
}

func NewLeaderBoardRepository(db *gorm.DB, logger *providers.ConsoleLogger) *LeaderboardRepository {
	return &LeaderboardRepository{
		db:     db,
		logger: logger,
	}
}

func (r *LeaderboardRepository) SubmitScore(ctx context.Context, userID int64, score int64, gameMode string) (time.Time, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return time.Time{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ Check user exists
	var exists bool
	if err := tx.Raw(
		`SELECT EXISTS (SELECT 1 FROM gaming.users WHERE id = ?)`,
		userID,
	).Scan(&exists).Error; err != nil {
		tx.Rollback()
		return time.Time{}, err
	}
	if !exists {
		tx.Rollback()
		return time.Time{}, errors.New(constants.ErrUserNotFound)
	}

	// 2️⃣ Insert game session
	now := time.Now().UTC()
	if err := tx.Exec(`
		INSERT INTO gaming.game_sessions (user_id, score, game_mode, timestamp)
		VALUES (?, ?, ?, ?)
	`, userID, score, gameMode, now).Error; err != nil {
		tx.Rollback()
		return time.Time{}, err
	}

	// 3️⃣ Update total_score in leaderboard
	// Using UPSERT: insert if not exists, otherwise increment
	if err := tx.Exec(`
		INSERT INTO gaming.leaderboard (user_id, total_score)
		VALUES (?, ?)
		ON CONFLICT (user_id)
		DO UPDATE SET total_score = leaderboard.total_score + EXCLUDED.total_score
	`, userID, score).Error; err != nil {
		tx.Rollback()
		return time.Time{}, err
	}

	// 4️⃣ Optional: recalculate ranks
	if err := tx.Exec(`
		WITH ranked AS (
			SELECT id, RANK() OVER (ORDER BY total_score DESC) AS rank
			FROM gaming.leaderboard
		)
		UPDATE gaming.leaderboard l
		SET rank = r.rank
		FROM ranked r
		WHERE l.id = r.id
	`).Error; err != nil {
		tx.Rollback()
		return time.Time{}, err
	}

	return now, tx.Commit().Error
}

func (r *LeaderboardRepository) GetTopPlayers(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	var entries []LeaderboardEntry

	err := r.db.Debug().WithContext(ctx).
		Table("gaming.leaderboard").
		Select("user_id, total_score").
		Order("total_score DESC").
		Limit(limit).
		Find(&entries).Error

	if err != nil {
		return nil, err
	}

	return entries, nil
}
