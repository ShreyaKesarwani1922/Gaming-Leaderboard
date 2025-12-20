// core/core.go
package core

import (
	"context"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/constants"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/model"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/repository"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
)

type LeaderboardCore struct {
	repo   *repository.LeaderboardRepository
	logger *providers.ConsoleLogger
}

type ILeaderboardCore interface {
	SubmitScore(ctx context.Context, req *model.SubmitScoreRequest) (*model.SubmitScoreResponse, error)
}

func NewLeaderboardCore(repo *repository.LeaderboardRepository, logger *providers.ConsoleLogger) *LeaderboardCore {
	return &LeaderboardCore{
		repo:   repo,
		logger: logger,
	}
}

func (c *LeaderboardCore) SubmitScore(ctx context.Context, req *model.SubmitScoreRequest) (*model.SubmitScoreResponse, error) {
	if req.Score < 0 {
		return &model.SubmitScoreResponse{
			Success: false,
			Error:   "Invalid score",
			Code:    constants.ErrInvalidScore,
		}, nil
	}

	timestamp, err := c.repo.SubmitScore(
		ctx,
		req.UserID,
		req.Score,
		req.GameMode,
	)
	if err != nil {
		return nil, err
	}

	return &model.SubmitScoreResponse{
		Success: true,
		Message: "Score submitted successfully",
		Data: &model.ScoreData{
			UserID:    req.UserID,
			Score:     req.Score,
			Timestamp: timestamp,
		},
	}, nil
}
