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

	// create and connect to database
	dbClient, err := db.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		if err := dbClient.Close(ctx); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	s := api.NewServer(ctx, os.Getenv("PORT"), dbClient)
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
