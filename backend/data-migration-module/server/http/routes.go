package http

import (
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
)

func (h *MigrationHandler) RegisterRoutes(router *mux.Router) {
	_, handler := newrelic.WrapHandle(h.newRelic, "/api/migrate/populate", http.HandlerFunc(h.PopulateData))
	router.Handle("/api/migrate/populate", handler).Methods("POST")
}
