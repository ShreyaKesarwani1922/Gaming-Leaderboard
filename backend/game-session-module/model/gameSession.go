package model

import (
	userModel "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/user-module/model"
	"time"
)

// GameSession represents the game_sessions table
type GameSession struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Score     int       `gorm:"not null" json:"score"`
	GameMode  string    `gorm:"type:varchar(50);not null" json:"game_mode"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`

	User userModel.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
