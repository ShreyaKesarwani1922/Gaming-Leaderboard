package main

import (
	"github.com/ShreyaKesarwani1922/Gaming-Leaderboard/user-module/server/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	//core := NewUserCore() // your implementation
	//metrics := nil        // your metrics implementation

	//httpHandler := http.Init(router, core, metrics)
	httpHandler := http.Init(router)

	http.ListenAndServe(":8080", router)
}
