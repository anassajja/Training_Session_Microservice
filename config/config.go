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

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		log.Fatal("MONGO_URI environment variable is required")
	}

	serverPortStr := os.Getenv("SERVER_PORT")
	if serverPortStr == "" {
		serverPortStr = "8080"
	}

	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		log.Fatalf("Invalid SERVER_PORT value: %v", err)
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		databaseName = "test"
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is required") // Log an error message if the JWT_SECRET_KEY is not set
	}

	log.Println("Mongo URI:", mongoURI)
	log.Println("Database Name:", databaseName)
	log.Println("Server Port:", serverPort)
	log.Println("JWT Secret Key:", jwtSecretKey)

	return &Config{
		ServerPort:   serverPort,
		MongoURI:     mongoURI,
		DatabaseName: databaseName,
		JwtSecretKey: jwtSecretKey,
	}
}
