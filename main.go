package main

import (
	"log"
	"os"

	"users-ms/data"
	"users-ms/middlewares"
	"users-ms/server"
)

func main() {
	store, err := data.NewStore(os.Getenv("DB_PROD"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Panic(err)
	}
	store.Populate()

	logger := middlewares.NewLogger()

	s := *server.NewServer(store, logger)

	s.Serve(os.Getenv("ADDR"))
}
