package model

import "gorm.io/gorm"

type GameSession struct {
	gorm.Model
	UserID    uint   `gorm:"column:user_id;not null;index"`
	Score     int    `gorm:"column:score;not null"`
	GameMode  string `gorm:"column:game_mode;type:varchar(50);not null"`
	Timestamp int64  `gorm:"column:timestamp;autoCreateTime"`
}
