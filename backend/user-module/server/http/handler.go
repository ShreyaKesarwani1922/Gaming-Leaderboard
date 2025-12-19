package http

import (
	"net/http"

	core "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/backend/user-module/core"
	"github.com/gorilla/mux"
)

// UserHttpExtension holds dependencies for the User HTTP module
type UserHttpExtension struct {
	Server  interface{} // replace with your server type
	Router  *mux.Router
	Metrics interface{} // replace with your metrics type
	Core    core.Core   // replace with your core service interface
}

// NewUserHttpExtension creates a new UserHttpExtension instance
func NewUserHttpExtension(router *mux.Router, core *core.Core) *UserHttpExtension {
	return &UserHttpExtension{
		Router: router,
		Core:   *core,
	}
}

// Init initializes the handler (to be called from main.go)
func (he *UserHttpExtension) Init() {
	he.Router.HandleFunc("/ping", he.Ping).Methods("GET")
}

// Ping handles the ping endpoint
func (he *UserHttpExtension) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "message": "pong"}`))
}
