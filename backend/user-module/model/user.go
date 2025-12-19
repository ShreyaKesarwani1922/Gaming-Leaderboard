package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(255);unique;not null"`
	JoinDate int64  `gorm:"column:join_date;autoCreateTime"`
}
