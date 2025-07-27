package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/obrown4/credit-stack/internal/db"
	"github.com/obrown4/credit-stack/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// User represents a user in the system
type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// Session represents a user session
type Session struct {
	Username     string    `bson:"username"`
	SessionToken string    `bson:"session_token"`
	CSRFToken    string    `bson:"csrf_token"`
	CreatedAt    time.Time `bson:"created_at"`
	ExpiresAt    time.Time `bson:"expires_at"`
}

// LoginResult contains the tokens generated during login
type LoginResult struct {
	SessionToken string
	CSRFToken    string
	Username     string
}

// RegisterUser handles the business logic for user registration
func RegisterUser(ctx context.Context, client *db.Client, username, password string) error {
	// Validate input
	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	if len(username) < 8 || len(password) < 8 {
		return fmt.Errorf("username and password must be at least 8 characters long")
	}

	users := client.Collection("creditStack", "users")
	log.Printf("Connected to users collection")

	// check if user already exists
	count, err := users.CountDocuments(ctx, bson.D{{Key: "username", Value: username}})
	if err != nil {
		return fmt.Errorf("failed to check if user already exists: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := User{
		Username: username,
		Password: hashedPassword,
	}

	_, err = users.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("User registered successfully")

	return nil
}

// LoginUser handles the business logic for user login
func LoginUser(ctx context.Context, client *db.Client, username, password string) (*LoginResult, error) {
	// Validate input
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// Check credentials against database
	users := client.Collection("creditStack", "users")
	count, err := users.CountDocuments(ctx, bson.D{{Key: "username", Value: username}})
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if count != 1 {
		return nil, fmt.Errorf("user does not exist")
	}

	user, err := getUser(ctx, client, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !utils.ComparePassword(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	// Generate tokens
	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

	// Store session in separate collection
	sessions := client.Collection("creditStack", "sessions")
	session := Session{
		Username:     username,
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	_, err = sessions.InsertOne(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	log.Printf("User logged in successfully")

	return &LoginResult{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		Username:     username,
	}, nil
}

// AuthorizeUser validates a user's session
func AuthorizeUser(ctx context.Context, client *db.Client, username, sessionToken, csrfToken string) error {
	if username == "" || sessionToken == "" || csrfToken == "" {
		return fmt.Errorf("username, sessionToken, and csrfToken are required")
	}

	// Check if user exists
	_, err := getUser(ctx, client, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Validate session
	sessions := client.Collection("creditStack", "sessions")
	sessionResult := sessions.FindOne(ctx, bson.D{
		{Key: "username", Value: username},
		{Key: "session_token", Value: sessionToken},
		{Key: "csrf_token", Value: csrfToken},
		{Key: "expires_at", Value: bson.M{"$gt": time.Now()}},
	})

	if sessionResult.Err() != nil {
		return fmt.Errorf("invalid or expired session")
	}

	var session Session
	err = sessionResult.Decode(&session)
	if err != nil {
		return fmt.Errorf("failed to decode session: %w", err)
	}

	return nil
}

// LogoutUser handles the business logic for user logout
func LogoutUser(ctx context.Context, client *db.Client, username, sessionToken string) error {
	if username == "" || sessionToken == "" {
		return fmt.Errorf("username and sessionToken are required")
	}

	// Remove session from database
	sessions := client.Collection("creditStack", "sessions")
	_, err := sessions.DeleteOne(ctx, bson.D{
		{Key: "username", Value: username},
		{Key: "session_token", Value: sessionToken},
	})

	if err != nil {
		return fmt.Errorf("failed to remove session: %w", err)
	}

	log.Printf("User logged out successfully")
	return nil
}

// getUser retrieves a user from the database
func getUser(ctx context.Context, client *db.Client, username string) (*User, error) {
	users := client.Collection("creditStack", "users")
	userResult := users.FindOne(ctx, bson.D{{Key: "username", Value: username}})
	if userResult.Err() != nil {
		return nil, fmt.Errorf("failed to find user: %w", userResult.Err())
	}

	var user User
	err := userResult.Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}
	return &user, nil
}

// PrintMessage handles the business logic for printing messages
func PrintMessage(ctx context.Context, msg string) {
	fmt.Println("POST request received!")
	fmt.Println("Message:", msg)
}
