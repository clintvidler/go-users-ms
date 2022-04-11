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

	ar.HandleFunc("/logout", uh.Logout).Methods(http.MethodPost)

	ar.HandleFunc("/", uh.CurrentUser).Methods(http.MethodGet)

	ar.HandleFunc("/update-info", uh.UpdateInfo).Methods(http.MethodPost)

	ar.HandleFunc("/update-password", uh.UpdatePassword).Methods(http.MethodPost)

	s.router = r
}
