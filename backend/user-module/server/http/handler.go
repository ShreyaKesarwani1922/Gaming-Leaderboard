package http

import (
	"encoding/json"
	core "github.com/ShreyaKesarwani1922/Gaming-Leaderboard/user-module/core"
	"net/http"

	"github.com/gorilla/mux"
)

// UserHttpExtension holds dependencies for the User HTTP module
type UserHttpExtension struct {
	Server  interface{} // replace with your server type
	Router  *mux.Router
	Metrics interface{} // replace with your metrics type
	Core    core.Core   // replace with your core service interface
}

// Init initializes the handler (to be called from main.go)
func (he *UserHttpExtension) Init() {
	he := &UserHttpExtension{
		Router:  router,
		Metrics: metrics,
		Core:    core,
	}
	RegisterRoutes(he)
	return he
}

// Example handler function
func (he *UserHttpExtension) CreateUserHandler() func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, rq *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		switch rq.Method {
		case http.MethodPost:
			var req CreateUserRequest
			if err := json.NewDecoder(rq.Body).Decode(&req); err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}

			// call core logic (placeholder)
			res, err := he.Core.(UserCore).CreateUser(req)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(rw).Encode(res)
		}
	}
}
