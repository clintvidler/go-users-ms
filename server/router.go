package server

import (
	"fmt"
	"net/http"
	"users-ms/middlewares"

	"github.com/gorilla/mux"
)

func (s *Server) setupRoutes() {
	r := mux.NewRouter().StrictSlash(true)

	// unauthenticated routes

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "501: Not Implemented")
	}).Methods(http.MethodPost)

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "501: Not Implemented")
	}).Methods(http.MethodPost)

	// authenticated routes

	ar := r.NewRoute().Subrouter()
	ar.Use(middlewares.IsAuthenticated)

	ar.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)

	ar.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodPost)

	ar.HandleFunc("/update-info", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodPost)

	ar.HandleFunc("/update-password", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodPost)

	s.Router = r
}
