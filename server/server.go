package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"users-ms/middlewares"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Logger *middlewares.Logger
}

func NewServer(logger middlewares.Logger) (s *Server) {
	s = &Server{Logger: &logger}
	return
}

func (s *Server) Serve(addr string) {
	s.setupRoutes()

	srv := &http.Server{
		Handler:      middlewares.LogRequestResponse(s.Router, *s.Logger),
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
