package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/obrown4/credit-stack/internal/db"
	"github.com/obrown4/credit-stack/internal/handlers"
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
	log.Printf("Connected to the database successfully")

	// set up handlers
	http.HandleFunc("/print", handlers.PrintMsg)

	// start network server
	http.ListenAndServe(":8080", nil)
}
