package main

import (
	"context"

	"loyalty-service/internal/handlers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := handlers.NewServer()
	server.Prepare(ctx)
	server.Run()
}
