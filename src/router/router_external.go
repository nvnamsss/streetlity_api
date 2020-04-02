package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func external(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("External"))
}

func Handle(router *mux.Router) {
	s := router.PathPrefix("/external").Subrouter()
	s.HandleFunc("/", external).
		Methods("GET", "POST")

	HandleFuel(router)
	HandlePing(router)
}