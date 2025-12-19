package model

import "gorm.io/gorm"

type LeaderBoard struct {
	gorm.Model
	UserID     uint `gorm:"column:user_id;not null;index"`
	TotalScore int  `gorm:"column:total_score;not null"`
	Rank       int  `gorm:"column:rank"`
}
