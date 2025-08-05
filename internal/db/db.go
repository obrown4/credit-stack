package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client wraps the MongoDB client with connection management
type Client struct {
	client *mongo.Client
}

// NewClient creates and connects a new MongoDB client
func NewClient(ctx context.Context) (*Client, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI environment variable is required")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Successfully connected to MongoDB")
	return &Client{client: client}, nil
}

// GetClient returns the underlying MongoDB client
func (c *Client) GetClient() *mongo.Client {
	return c.client
}

// Close disconnects the MongoDB client
func (c *Client) Close(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Disconnect(ctx)
}

// Database returns a database instance
func (c *Client) Database(name string) *mongo.Database {
	return c.client.Database(name)
}

// Collection returns a collection instance
func (c *Client) Collection(dbName, collName string) *mongo.Collection {
	return c.client.Database(dbName).Collection(collName)
}

// IsPresent checks if a document exists in a collection
func IsPresent(collection *mongo.Collection, filter interface{}) (bool, error) {
	err := collection.FindOne(context.TODO(), filter).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
