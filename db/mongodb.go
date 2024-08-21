package db

import (
	"context"
	"log"
	"training_session/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client // Client

func Connect() { // Connect to MongoDB
	cfg := config.LoadConfig() // Load configuration

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)   // Set the MongoDB URI
	var err error                                              // Declare an error variable
	Client, err = mongo.Connect(context.TODO(), clientOptions) // Connect to MongoDB
	if err != nil {                                            // Check if there is an error
		log.Fatal("Failed to connect to MongoDB:", err) // Log an error if the connection fails
	}
	err = Client.Ping(context.TODO(), nil) // Ping MongoDB
	if err != nil {                        // Check if there is an error
		log.Fatal("Failed to ping MongoDB:", err) // Log an error if the ping fails
	}
	log.Println("Connected to MongoDB!") // Log a message if the connection is successful
}
