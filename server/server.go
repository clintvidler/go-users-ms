package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"users-ms/data"
	"users-ms/middlewares"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	router    *mux.Router
	datastore *data.Store
	logger    *middlewares.Logger
}

func NewServer(ds *data.Store, logger *middlewares.Logger) (s *Server) {
	s = &Server{datastore: ds, logger: logger}
	return
}

func (s *Server) Serve(addr string) {
	s.setupRoutes()

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	headers := []string{"Content-Type"}
	origins := []string{"http://localhost:9090"} // TODO: review allowed origins 	// origins := []string{"*"}

	cors := handlers.CORS(
		handlers.AllowedMethods(methods),
		handlers.AllowedHeaders(headers),
		handlers.AllowedOrigins(origins),
		handlers.AllowCredentials(),
	)

	srv := &http.Server{
		Handler:      cors(middlewares.LogRequestResponse(s.router, *s.logger)),
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		log.Println("Starting server on port 9090")

		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm/interupt and gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting a time for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(ctx)
}
