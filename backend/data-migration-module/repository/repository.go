package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type IMigrationRepository interface {
	BulkInsertUsers(tx *gorm.DB, limit int) error
	BulkInsertGameSessions(tx *gorm.DB, limit int) error
	UpdateLeaderboard(tx *gorm.DB) error
	GetMaxUserID(tx *gorm.DB) (int64, error)
}

type MigrationRepository struct {
	db *gorm.DB
}

func NewMigrationRepository(db *gorm.DB) *MigrationRepository {
	return &MigrationRepository{db: db}
}

func (r *MigrationRepository) BulkInsertUsers(tx *gorm.DB, limit int) error {
	sql := `
		INSERT INTO gaming.users (username)
		SELECT 'user_' || (SELECT COALESCE(MAX(id),0) FROM gaming.users) + generate_series(1, ?);
	`
	return tx.Exec(sql, limit).Error
}

func (r *MigrationRepository) BulkInsertGameSessions(tx *gorm.DB, limit int) error {
	sql := `
		INSERT INTO gaming.game_sessions (user_id, score, game_mode, timestamp)
		SELECT
			u.id,
			floor(random() * 10000 + 1)::int,
			CASE WHEN random() > 0.5 THEN 'solo' ELSE 'team' END,
			NOW() - INTERVAL '1 day' * floor(random() * 365)
		FROM (
			SELECT id
			FROM gaming.users
			ORDER BY random()
			LIMIT ?
		) u;
	`
	return tx.Exec(sql, limit).Error
}

func (r *MigrationRepository) UpdateLeaderboard(tx *gorm.DB) error {
	sql := `
		INSERT INTO gaming.leaderboard (user_id, total_score, rank)
		SELECT
			user_id,
			SUM(score) AS total_score,
			RANK() OVER (ORDER BY SUM(score) DESC)
		FROM gaming.game_sessions
		GROUP BY user_id
		ON CONFLICT (user_id) DO UPDATE
		SET
			total_score = EXCLUDED.total_score,
			rank = EXCLUDED.rank;
	`
	return tx.Exec(sql).Error
}

func (r *MigrationRepository) GetMaxUserID(tx *gorm.DB) (int64, error) {
	var maxUserID int64
	err := tx.Raw("SELECT COALESCE(MAX(id), 0) FROM gaming.users").Scan(&maxUserID).Error
	if err != nil {
		return 0, err
	}
	if maxUserID == 0 {
		return 0, fmt.Errorf("no users found in users table")
	}
	return maxUserID, nil
}
