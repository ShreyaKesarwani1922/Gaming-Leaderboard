package model

import "time"

type SubmitScoreRequest struct {
	UserID   int64  `json:"user_id" validate:"required"`
	Score    int64  `json:"score" validate:"required,min=0"`
	GameMode string `json:"game_mode" validate:"required,oneof=solo team"` // Add more modes as needed
}

type SubmitScoreResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    *ScoreData `json:"data,omitempty"`
	Error   string     `json:"error,omitempty"`
	Code    string     `json:"code,omitempty"`
}

type ScoreData struct {
	UserID    int64     `json:"user_id"`
	Score     int64     `json:"score"`
	Timestamp time.Time `json:"timestamp"`
}
