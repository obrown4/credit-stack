package auth

import (
	"fmt"

	"github.com/obrown4/credit-stack/internal/db"
)

// User represents a user in the system
type User struct {
	Username string
	Password string
}

// RegisterUser handles the business logic for user registration
func RegisterUser(username, password string) error {
	// Validate input
	if len(username) < 8 || len(password) < 8 {
		return fmt.Errorf("username and password must be at least 8 characters long")
	}

	// TODO: Add database integration
	users := db.Client.Database("creditStack").Collection("users")

	return nil
}

// LoginUser handles the business logic for user login
func LoginUser(username, password string) error {
	// Validate input
	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	// TODO: Add authentication logic
	// Check credentials against database

	return nil
}

// LogoutUser handles the business logic for user logout
func LogoutUser() error {
	// TODO: Add logout logic (clear session, etc.)
	return nil
}

// PrintMessage handles the business logic for printing messages
func PrintMessage(msg string) {
	fmt.Println("POST request received!")
	fmt.Println("Message:", msg)
}
