package main

import (
	"log"
	"training_session/db"
	"training_session/pkg/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	db.Connect()

	// Check if the client is initialized
	if db.Client == nil {
		log.Fatal("MongoDB client is not initialized")
	}

	// Set up routes and start the server
	r := gin.Default()
	routes.SetupRoutes(r)

	// Run the server
	log.Println("Starting server on :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
