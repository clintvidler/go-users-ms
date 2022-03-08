package main

import (
	"os"

	"users-ms/server"
)

func main() {
	logger := *server.NewLogger()

	s := *server.NewServer(logger)

	s.Serve(os.Getenv("ADDR"))
}
