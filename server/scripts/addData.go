package main

import (
	//"context"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/obrown4/credit-stack/server/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RewardsCategory struct {
	Category           string   `json:"category"`           // e.g., "Dining", "Travel"
	RewardRate         float64  `json:"rewardRate"`         // e.g., 3.0 for 3%
	Exclusions         []string `json:"exclusions"`         // e.g., ["Target", "Walmart"]
	RotatingCategories []string `json:"rotatingCategories"` // ["Groceries"]
}

type CreditCard struct {
	Name              string            `json:"name"`
	Issuer            string            `json:"issuer"`
	RewardsCategories []RewardsCategory `json:"rewardsCategories"`
}

func main() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(env.GetConnectionString()).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	log.Printf("Connected to Atlas\n")

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("rewardsDB").Collection("creditCards")

	addCards(collection)
}

func addCards(col *mongo.Collection) {
	var creditCards []CreditCard

	data, err := os.ReadFile("creditCards.json")
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	err = json.Unmarshal(data, &creditCards)
	if err != nil {
		fmt.Printf("Error parsing file: %s\n", err)
		return
	}

	for _, card := range creditCards {
		// Create filter based on card name (or whatever field you want to use as unique identifier)
		filter := bson.M{"name": card.Name}

		// Create update document
		update := bson.M{"$set": card}

		// Set upsert option to true
		opts := options.Update().SetUpsert(true)

		// Perform upsert
		_, err := col.UpdateOne(
			context.TODO(),
			filter,
			update,
			opts,
		)

		if err != nil {
			log.Printf("%s\n", err)
			return
		}
	}
}
