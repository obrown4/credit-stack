package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/obrown4/credit-stack/api"
	"github.com/obrown4/credit-stack/internal/db"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to db
	if err := db.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Disconnect(ctx)

	log.Printf("Connected to the database successfully")

	s := api.NewServer(ctx, ":8080")
	var eg errgroup.Group

	// start network server
	eg.Go(func() error {
		return s.Run()
	})

	<-ctx.Done()

	if err := s.Close(); err != nil {
		log.Fatalf("Failed to close the server: %v", err)
	}

	if err := eg.Wait(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Printf("Server closed")
}
