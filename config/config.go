package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   int    // Server port number
	MongoURI     string // MongoDB URI
	DatabaseName string // Database name
	JwtSecretKey string // JWT secret key
}

func LoadConfig() *Config { // LoadConfig function to load the configuration
	if err := godotenv.Load(); err != nil { // Load the .env file
		log.Fatal("Error loading .env file") // Log an error message if the .env file cannot be loaded
	}

	mongoURI, exists := os.LookupEnv("MONGO_URI") // Get the MongoDB URI from the environment
	if !exists {                                  // Check if the MONGO_URI environment variable is set
		log.Fatal("MONGO_URI environment variable is required") // Log an error message if the MONGO_URI is not set
	}

	serverPortStr := os.Getenv("SERVER_PORT") // Get the server port from the environment
	if serverPortStr == "" {                  // Check if the SERVER_PORT environment variable is not set
		serverPortStr = "8080" // Set the default server port to 8080
	}

	serverPort, err := strconv.Atoi(serverPortStr) // Convert the server port to an integer
	if err != nil {                                // Check if there is an error
		log.Fatalf("Invalid SERVER_PORT value: %v", err) // Log an error message if the SERVER_PORT value is invalid
	}

	databaseName := os.Getenv("DATABASE_NAME") // Get the database name from the environment
	if databaseName == "" {                    // Check if the DATABASE_NAME environment variable is not set
		databaseName = "test" // Set the default database name to "test"
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY") // Get the JWT secret key from the environment
	if jwtSecretKey == "" {                     // Check if the JWT_SECRET_KEY environment variable is not set
		log.Fatal("JWT_SECRET_KEY environment variable is required") // Log an error message if the JWT_SECRET_KEY is not set
	}

	log.Println("Mongo URI:", mongoURI)          // Log the MongoDB URI
	log.Println("Database Name:", databaseName)  // Log the database name
	log.Println("Server Port:", serverPort)      // Log the server port
	log.Println("JWT Secret Key:", jwtSecretKey) // Log the JWT secret key

	return &Config{ // Return the configuration
		ServerPort:   serverPort,   // Set the server port
		MongoURI:     mongoURI,     // Set the MongoDB URI
		DatabaseName: databaseName, // Set the database name
		JwtSecretKey: jwtSecretKey, // Set the JWT secret key
	}
}
