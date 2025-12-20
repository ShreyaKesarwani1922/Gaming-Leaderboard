package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func (h *LeaderboardHandler) RegisterRoutes(router *mux.Router) {
	// Submit score endpoint
	_, submitHandler := newrelic.WrapHandle(h.newrelic, "api/leaderboard/submit", http.HandlerFunc(h.SubmitScore))
	router.Handle("/api/leaderboard/submit", submitHandler).Methods(http.MethodPost)

	// Get top players endpoint
	_, topPlayersHandler := newrelic.WrapHandle(h.newrelic, "api/leaderboard/top", http.HandlerFunc(h.GetTopPlayers))
	router.Handle("/api/leaderboard/top", topPlayersHandler).Methods(http.MethodGet)

	// Get player rank endpoint
	_, playerRankHandler := newrelic.WrapHandle(h.newrelic, "api/leaderboard/rank/{user_id}", http.HandlerFunc(h.GetPlayerRank))
	router.Handle("/api/leaderboard/rank/{user_id}", playerRankHandler).Methods(http.MethodGet)

	// Get leaderboard stream endpoint
	_, leaderboardStreamHandler := newrelic.WrapHandle(h.newrelic, "api/leaderboard/stream", http.HandlerFunc(h.StreamLeaderboard))
	router.Handle("/api/leaderboard/stream", leaderboardStreamHandler).Methods(http.MethodGet)

	
}
