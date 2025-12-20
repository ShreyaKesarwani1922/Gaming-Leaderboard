package http

import (
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
)

func (h *LeaderboardHandler) RegisterRoutes(router *mux.Router) {
	_, handler := newrelic.WrapHandle(h.newrelic, "api/leaderboard/submit", http.HandlerFunc(h.SubmitScore))
	router.Handle("/api/leaderboard/submit", handler).Methods(http.MethodPost)
}
