package model

import "time"

type User struct {
	ID               uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Username         string        `gorm:"type:varchar(255);unique;not null" json:"username"`
	JoinDate         time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"join_date"`
	Sessions         []GameSession `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	LeaderboardEntry Leaderboard   `gorm:"foreignKey:UserID" json:"leaderboard_entry,omitempty"`
}
