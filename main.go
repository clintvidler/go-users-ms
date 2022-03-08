package main

import (
	"os"

	"users-ms/server"
)

func main() {
	s := server.NewServer()

	s.Serve(os.Getenv("ADDR"))
}
