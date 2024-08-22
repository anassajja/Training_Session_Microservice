package db

import (
	"context"
	"log"
	"time"
	"training_session/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client   *mongo.Client   // Global MongoDB client variable
	Database *mongo.Database // Global MongoDB database variable
)

// Connect initializes the MongoDB connection using the provided configuration.
func Connect(cfg *config.Config) (*mongo.Database, error) {
	// Create a context with a timeout for the MongoDB connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Create a context with a timeout
	defer cancel()                                                           // Defer the cancel function to release resources after the function returns.

	// Set the client options using the MongoDB URI from the configuration.
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	// Connect to MongoDB using the client options.
	client, err := mongo.Connect(ctx, clientOptions) // Connect to MongoDB
	if err != nil {                                  // Check if there is an error
		return nil, err // Return the error
	}

	// Ping the MongoDB server to verify the connection.
	err = client.Ping(ctx, readpref.Primary()) // Ping the MongoDB server
	if err != nil {                            // Check if there is an error
		return nil, err // Return the error
	}

	// Set the global client and database variables.
	Client = client                              // Set the global client variable
	Database = client.Database(cfg.DatabaseName) // Set the global database variable

	// Log a success message and return the connected database.
	log.Println("Connected to MongoDB!") // Log a success message
	return Database, nil                 // Return the connected database
}
