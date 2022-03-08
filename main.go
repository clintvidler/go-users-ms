package main

import (
	"os"

	"users-ms/middlewares"
	"users-ms/server"
)

func main() {
	logger := *middlewares.NewLogger()

	s := *server.NewServer(logger)

	s.Serve(os.Getenv("ADDR"))
}
