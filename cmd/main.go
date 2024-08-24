package main

import (
	"fmt"
	"log"
	"training_session/config"
	"training_session/db"
	"training_session/pkg/controllers"
	"training_session/pkg/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig() // Load the configuration

	// Connect to MongoDB
	database, err := db.Connect(cfg) // Connect to MongoDB
	if err != nil {                  // Check if there is an error
		log.Fatalf("Failed to connect to MongoDB: %v", err) // Log the error message
	}

	// Log message indicating that the database has been initialized successfully
	log.Println("Database initialized successfully") // Add this log message

	// Initialize controllers with the database connection
	controllers.InitializeSession(database) // Initialize the controllers
	controllers.InitializeUser(database)    // Initialize the controllers
	controllers.InitializeInvitation(database)
	controllers.InitializeNotification(database)
	controllers.InitializeFeedbackController(database)
	controllers.InitializePitchBooking(database)

	// Set up routes and start the server
	r := gin.Default()    // Create a new Gin router
	routes.SetupRoutes(r) // Set up the routes

	// Run the server
	log.Printf("Starting server on :%d", cfg.ServerPort)              // Log the server port
	if err := r.Run(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil { // Run the server
		log.Fatal(err) // Log any errors
	}
}
