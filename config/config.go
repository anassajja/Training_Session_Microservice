package config 

import (
	"log"
	"os"
)

// Config holds the application configuration
type Config struct { // Define a Config struct
	MongoURI string // Define a MongoURI field
}

func LoadConfig() *Config { // Load configuration
	mongoURI, exists := os.LookupEnv("MONGO_URI") // Get the value of the MONGO_URI environment variable
	if !exists { 					// Check if the environment variable exists
		log.Fatal("MONGO_URI environment variable is required") // Log an error if the environment variable is not found
	}

	log.Println("Mongo URI:", mongoURI) // Log the Mongo URI

	return &Config{ // Return a new Config instance
		MongoURI: mongoURI, // Set the MongoURI field
	}
}
