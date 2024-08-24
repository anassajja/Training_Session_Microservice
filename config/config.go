package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct { // Define a Config struct
	ServerPort   int    // Define a ServerPort field of type int
	MongoURI     string // Define a MongoURI field of type string
	DatabaseName string // Define a DatabaseName field of type string
}

func LoadConfig() *Config { // Load the configuration
	if err := godotenv.Load(); err != nil { // Load the .env file
		log.Fatal("Error loading .env file") // Log an error message
	}

	mongoURI, exists := os.LookupEnv("MONGO_URI") // Get the MONGO_URI environment variable
	if !exists {                                  // Check if the environment variable does not exist
		log.Fatal("MONGO_URI environment variable is required") // Log an error message
	}

	serverPortStr := os.Getenv("SERVER_PORT") // Get the SERVER_PORT environment variable
	if serverPortStr == "" {                  // Check if the environment variable is empty
		serverPortStr = "8080" // Set the default server port
	}

	serverPort, err := strconv.Atoi(serverPortStr) // Convert the server port to an integer
	if err != nil {
		log.Fatalf("Invalid SERVER_PORT value: %v", err) // Log an error message
	}

	databaseName := os.Getenv("DATABASE_NAME") // Get the DATABASE_NAME environment variable
	if databaseName == "" {                    // Check if the environment variable is empty
		databaseName = "test" // Set the default database name
	} // Set the default database name

	log.Println("Mongo URI:", mongoURI)         // Log the MongoDB URI
	log.Println("Database Name:", databaseName) // Log the database name
	log.Println("Server Port:", serverPort)     // Log the server port

	return &Config{
		ServerPort:   serverPort,
		MongoURI:     mongoURI,
		DatabaseName: databaseName,
		
	}
}
