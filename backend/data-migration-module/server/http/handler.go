package http

import (
	"encoding/json"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/core"
)

type MigrationHandler struct {
	core     *core.MigrationCore
	newRelic *newrelic.Application
}

func NewMigrationHandler(core *core.MigrationCore, newrelic *newrelic.Application) *MigrationHandler {
	return &MigrationHandler{
		core:     core,
		newRelic: newrelic,
	}
}

type PopulateGameSessionsRequest struct {
	SessionLimit int `json:"session_limit"`
}

type PopulateDataRequest struct {
	UserLimit    int `json:"user_limit"`
	SessionLimit int `json:"session_limit"`
}

type PopulateDataResponse struct {
	Message string `json:"message"`
}

func (h *MigrationHandler) PopulateData(w http.ResponseWriter, r *http.Request) {
	var req PopulateDataRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserLimit <= 0 || req.SessionLimit <= 0 {
		http.Error(w, "user_limit and session_limit must be greater than 0", http.StatusBadRequest)
		return
	}

	if req.UserLimit > 100_000 || req.SessionLimit > 500_000 {
		http.Error(w, "Limits too high for single request", http.StatusBadRequest)
		return
	}

	if err := h.core.PopulateSampleData(
		r.Context(),
		req.UserLimit,
		req.SessionLimit,
	); err != nil {
		http.Error(
			w,
			"Failed to populate data: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(PopulateDataResponse{
		Message: "Data population completed successfully",
	})
}
