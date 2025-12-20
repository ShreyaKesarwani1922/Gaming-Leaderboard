package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/constants"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/core"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/leader-board-module/model"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
	"go.uber.org/zap"
)

type LeaderboardHandler struct {
	core     *core.LeaderboardCore
	logger   *providers.ConsoleLogger
	newrelic *newrelic.Application
}

func NewLeaderboardHandler(core *core.LeaderboardCore, logger *providers.ConsoleLogger, newrelic *newrelic.Application) *LeaderboardHandler {
	return &LeaderboardHandler{
		core:     core,
		logger:   logger,
		newrelic: newrelic,
	}
}

func (h *LeaderboardHandler) SubmitScore(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()

	var req model.SubmitScoreRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.respondWithError(
			w,
			http.StatusBadRequest,
			"Invalid request payload",
			constants.ErrInvalidScore,
		)
		return
	}

	// Basic validation (keep minimal)
	if req.UserID <= 0 {
		h.respondWithError(
			w,
			http.StatusBadRequest,
			"Invalid user_id",
			constants.ErrUserNotFound,
		)
		return
	}

	resp, err := h.core.SubmitScore(ctx, &req)
	if err != nil {
		h.logger.Error(
			"SubmitScore failed",
			zap.Int64("user_id", req.UserID),
			zap.Error(err),
		)

		h.respondWithError(
			w,
			http.StatusInternalServerError,
			"Internal server error",
			constants.ErrInternalServer,
		)
		return
	}

	if !resp.Success {
		status := http.StatusBadRequest
		if resp.Code == constants.ErrUserNotFound {
			status = http.StatusNotFound
		}

		h.respondWithJSON(w, status, resp)
		return
	}

	h.respondWithJSON(w, http.StatusOK, resp)
}

func (h *LeaderboardHandler) respondWithJSON(
	w http.ResponseWriter,
	status int,
	payload interface{},
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *LeaderboardHandler) respondWithError(
	w http.ResponseWriter,
	status int,
	message string,
	code string,
) {
	h.respondWithJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
		"code":    code,
	})
}

func (h *LeaderboardHandler) GetTopPlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Default to top 10 players if limit is not specified or invalid
	limit := 10
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	resp, err := h.core.GetTopPlayers(ctx, limit)
	if err != nil {
		h.logger.Error(
			"GetTopPlayers failed",
			zap.Error(err),
		)

		h.respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch leaderboard",
			constants.ErrInternalServer,
		)
		return
	}

	h.respondWithJSON(w, http.StatusOK, resp)
}

func (h *LeaderboardHandler) GetPlayerRank(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		h.respondWithError(
			w,
			http.StatusBadRequest,
			"Invalid user ID",
			constants.ErrInvalidRequest,
		)
		return
	}

	resp, err := h.core.GetPlayerRank(ctx, userID)
	if err != nil {
		h.logger.Error(
			"GetPlayerRank failed",
			zap.Int64("user_id", userID),
			zap.Error(err),
		)

		h.respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to fetch player rank",
			constants.ErrInternalServer,
		)
		return
	}

	if !resp.Success {
		status := http.StatusNotFound
		if resp.Code == constants.ErrUserNotFound {
			status = http.StatusNotFound
		}

		h.respondWithJSON(w, status, resp)
		return
	}

	h.respondWithJSON(w, http.StatusOK, resp)
}
