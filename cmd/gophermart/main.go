package main

import (
	"loyalty-service/internal/handlers"
)

func main() {
	server := handlers.NewServer()
	server.Run()
}
