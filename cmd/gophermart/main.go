package main

import (
	"log"

	"loyalty-service/internal/handlers"
)

func main() {
	log.Println(111777)
	server := handlers.NewServer()
	server.Run()
}
