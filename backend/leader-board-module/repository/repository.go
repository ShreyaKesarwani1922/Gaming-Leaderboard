package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/constants"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	cacheTTL          = 5 * time.Minute
	maxRetries        = 3
	initialRetryDelay = 100 * time.Millisecond

	leaderboardVersionKey = "leaderboard:version"
	topPlayersCacheKey    = "leaderboard:top:%d:%d"    // version, limit
	playerRankCacheKey    = "leaderboard:player:%d:%d" // version, userID
)

type LeaderboardEntry struct {
	UserID     int64 `gorm:"column:user_id" json:"user_id"`
	TotalScore int64 `gorm:"column:total_score" json:"total_score"`
}

type PlayerRank struct {
	UserID int64 `gorm:"column:user_id" json:"user_id"`
	Rank   int   `json:"rank"`
	Score  int64 `gorm:"column:total_score" json:"score"`
}

type ILeaderboardRepository interface {
	SubmitScore(ctx context.Context, userID, score int64, gameMode string) (time.Time, error)
	GetTopPlayers(ctx context.Context, limit int) ([]LeaderboardEntry, error)
	GetPlayerRank(ctx context.Context, userID int64) (*PlayerRank, error)
}

type LeaderboardRepository struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *providers.ConsoleLogger
}

func NewLeaderBoardRepository(
	db *gorm.DB,
	redisClient *redis.Client,
	logger *providers.ConsoleLogger,
) *LeaderboardRepository {
	return &LeaderboardRepository{
		db:     db,
		redis:  redisClient,
		logger: logger,
	}
}

/* ============================
   Internal helpers
============================ */

func (r *LeaderboardRepository) leaderboardVersion(ctx context.Context) int64 {
	if r.redis == nil {
		return 1
	}

	version, err := r.redis.Get(ctx, leaderboardVersionKey).Int64()
	if err == redis.Nil {
		r.redis.Set(ctx, leaderboardVersionKey, 1, 0)
		return 1
	}
	if err != nil {
		r.logger.Warn("Failed to get leaderboard version", "error", err)
		return 1
	}
	return version
}

func (r *LeaderboardRepository) bumpLeaderboardVersion(ctx context.Context) {
	if r.redis == nil {
		return
	}
	if err := r.redis.Incr(ctx, leaderboardVersionKey).Err(); err != nil {
		r.logger.Warn("Failed to bump leaderboard version", "error", err)
	}
}

/* ============================
   Submit Score
============================ */

func (r *LeaderboardRepository) SubmitScore(
	ctx context.Context,
	userID int64,
	score int64,
	gameMode string,
) (time.Time, error) {

	var lastErr error
	now := time.Now().UTC()

	for attempt := 0; attempt < maxRetries; attempt++ {
		tx := r.db.WithContext(ctx).Begin()
		if tx.Error != nil {
			lastErr = tx.Error
			time.Sleep(initialRetryDelay * time.Duration(attempt+1))
			continue
		}

		// Check user exists (fast EXISTS)
		var exists bool
		if err := tx.Raw(
			`SELECT EXISTS (SELECT 1 FROM gaming.users WHERE id = ?)`,
			userID,
		).Scan(&exists).Error; err != nil || !exists {
			tx.Rollback()
			return time.Time{}, errors.New(constants.ErrUserNotFound)
		}

		// Insert game session
		if err := tx.Exec(`
			INSERT INTO gaming.game_sessions (user_id, score, game_mode, timestamp)
			VALUES (?, ?, ?, ?)
		`, userID, score, gameMode, now).Error; err != nil {
			tx.Rollback()
			lastErr = err
			time.Sleep(initialRetryDelay * time.Duration(attempt+1))
			continue
		}

		// Atomic upsert leaderboard score
		if err := tx.Exec(`
			INSERT INTO gaming.leaderboard (user_id, total_score)
			VALUES (?, ?)
			ON CONFLICT (user_id)
			DO UPDATE SET total_score = leaderboard.total_score + EXCLUDED.total_score
		`, userID, score).Error; err != nil {
			tx.Rollback()
			lastErr = err
			time.Sleep(initialRetryDelay * time.Duration(attempt+1))
			continue
		}

		if err := tx.Commit().Error; err != nil {
			lastErr = err
			time.Sleep(initialRetryDelay * time.Duration(attempt+1))
			continue
		}

		// Cache invalidation (O(1))
		r.bumpLeaderboardVersion(ctx)

		return now, nil
	}

	return time.Time{}, fmt.Errorf("submit score failed after retries: %w", lastErr)
}

/* ============================
   Get Top Players
============================ */

func (r *LeaderboardRepository) GetTopPlayers(
	ctx context.Context,
	limit int,
) ([]LeaderboardEntry, error) {

	if limit <= 0 {
		return nil, errors.New("limit must be positive")
	}

	version := r.leaderboardVersion(ctx)
	cacheKey := fmt.Sprintf(topPlayersCacheKey, version, limit)

	if r.redis != nil {
		if cached, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
			var entries []LeaderboardEntry
			if json.Unmarshal([]byte(cached), &entries) == nil {
				return entries, nil
			}
		}
	}

	var entries []LeaderboardEntry
	err := r.db.WithContext(ctx).
		Table("gaming.leaderboard").
		Select("user_id, total_score").
		Order("total_score DESC").
		Limit(limit).
		Find(&entries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch top players: %w", err)
	}

	if r.redis != nil {
		if data, err := json.Marshal(entries); err == nil {
			r.redis.Set(ctx, cacheKey, data, cacheTTL)
		}
	}

	return entries, nil
}

/* ============================
   Get Player Rank
============================ */

func (r *LeaderboardRepository) GetPlayerRank(
	ctx context.Context,
	userID int64,
) (*PlayerRank, error) {

	version := r.leaderboardVersion(ctx)
	cacheKey := fmt.Sprintf(playerRankCacheKey, version, userID)

	if r.redis != nil {
		if cached, err := r.redis.Get(ctx, cacheKey).Result(); err == nil {
			var rank PlayerRank
			if json.Unmarshal([]byte(cached), &rank) == nil {
				return &rank, nil
			}
		}
	}

	var rank PlayerRank
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			user_id,
			total_score,
			1 + (
				SELECT COUNT(*)
				FROM gaming.leaderboard
				WHERE total_score > lb.total_score
			) AS rank
		FROM gaming.leaderboard lb
		WHERE user_id = ?
	`, userID).Scan(&rank).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(constants.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player rank: %w", err)
	}

	if r.redis != nil {
		if data, err := json.Marshal(rank); err == nil {
			r.redis.Set(ctx, cacheKey, data, cacheTTL)
		}
	}

	return &rank, nil
}
