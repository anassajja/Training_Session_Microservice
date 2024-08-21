package config

import (
	"log"
	"os"
)

// Config holds the application configuration
type Config struct {
	MongoURI string
}

func LoadConfig() *Config { // Load configuration
	mongoURI, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		log.Fatal("MONGO_URI environment variable is required")
	}

	log.Println("Mongo URI:", mongoURI)

	return &Config{
		MongoURI: mongoURI,
	}
}
