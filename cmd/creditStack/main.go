package main

import (
	"log"

	"github.com/joho/godotenv"
	db "github.com/obrown4/credit-stack/internal/db"
)

func main() {
	// load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		panic(err)
	}

	log.Println("Database connection established successfully.")
}
