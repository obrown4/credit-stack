package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/obrown4/credit-stack/api"
	"github.com/obrown4/credit-stack/internal/db"
)

func main() {
	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to db
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Client.Disconnect(context.TODO())
	log.Printf("Connected to the database successfully")

	s := api.NewServer()
	api.SetRoutes(&s)

	// start network server
	http.ListenAndServe(":8080", &s)
}
