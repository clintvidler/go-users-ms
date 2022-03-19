package server

import (
	"net/http"
	"users-ms/handlers"

	"github.com/gorilla/mux"
)

func (s *Server) setupRoutes() {
	r := mux.NewRouter().StrictSlash(true)

	uh := handlers.NewUsersHandler(s.datastore, s.logger)

	th := handlers.NewTokensHandler(s.datastore, s.logger)

	// unauthenticated routes

	r.Handle("/login", uh.ValidateInfo(http.HandlerFunc(uh.Login))).Methods(http.MethodPost)

	r.Handle("/register", uh.ValidateInfo(http.HandlerFunc(uh.Register))).Methods(http.MethodPost)

	// authenticated routes

	ar := r.NewRoute().Subrouter()

	ar.Use(th.IsAuthenticated)

	ar.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)

	ar.HandleFunc("/logout", uh.Logout).Methods(http.MethodPost)

	ar.HandleFunc("/update-info", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodPost)

	ar.HandleFunc("/update-password", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodPost)

	s.router = r
}
