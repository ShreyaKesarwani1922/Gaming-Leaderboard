package http

import (
	"encoding/json"
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/providers"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"

	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/data-migration-module/core"
)

type MigrationHandler struct {
	core     *core.MigrationCore
	newRelic *newrelic.Application
	logger   *providers.ConsoleLogger
}

func NewMigrationHandler(core *core.MigrationCore, newrelic *newrelic.Application, logger *providers.ConsoleLogger) *MigrationHandler {
	return &MigrationHandler{
		core:     core,
		newRelic: newrelic,
		logger:   logger,
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
	h.logger.Info("PopulateData request received")

	var req PopulateDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info(
		"populate data request parsed: ",
		"user_limit: ", req.UserLimit,
		" session_limit: ", req.SessionLimit,
	)

	txn := newrelic.FromContext(r.Context())
	if txn != nil {
		txn.AddAttribute("user_limit", req.UserLimit)
		txn.AddAttribute("session_limit", req.SessionLimit)
	}

	if req.UserLimit <= 0 || req.SessionLimit <= 0 {
		h.logger.Warn(
			"invalid limits provided",
			"user_limit", req.UserLimit,
			"session_limit", req.SessionLimit,
		)
		http.Error(w, "user_limit and session_limit must be greater than 0", http.StatusBadRequest)
		return
	}

	if req.UserLimit > 100_000 || req.SessionLimit > 500_000 {
		h.logger.Warn(
			"limits too high for single request",
			"user_limit", req.UserLimit,
			"session_limit", req.SessionLimit,
		)
		http.Error(w, "Limits too high for single request", http.StatusBadRequest)
		return
	}

	h.logger.Info("starting sample data population")

	segment := newrelic.StartSegment(txn, "PopulateSampleData")
	err := h.core.PopulateSampleData(r.Context(), req.UserLimit, req.SessionLimit)
	segment.End()
	if err != nil {
		h.logger.Error(
			"failed to populate sample data",
			err,
			"user_limit", req.UserLimit,
			"session_limit", req.SessionLimit,
		)
		http.Error(
			w,
			"Failed to populate data: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	h.logger.Info(
		"sample data population completed successfully",
		"user_limit", req.UserLimit,
		"session_limit", req.SessionLimit,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(PopulateDataResponse{
		Message: "Data population completed successfully",
	})
}
