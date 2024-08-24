package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"training_session/pkg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var notificationCollection *mongo.Collection // Global notification collection

func InitializeNotification(database *mongo.Database) { // Initialize the notification controller
    if database == nil { // Check if the database is nil
        log.Println("Error: database is not initialized") // Log an error message
        return
    }
    // Initialize the collection
    notificationCollection = database.Collection("notifications") // Initialize the notification collection
    log.Println("notificationCollection initialized successfully") // Log a success message
}

// SendSessionNotification: Sends notifications to users about session changes or updates.
func SendSessionNotification(notification models.Notification) error { // Send a session notification
	// Ensure that notificationCollection is initialized
	if notificationCollection == nil { // Check if the notificationCollection is nil
		log.Println("Error: notificationCollection is not initialized") // Log an error message
		return mongo.ErrNilDocument                                     // Return a nil document error
	}

	// Insert the notification into the MongoDB collection
	_, err := notificationCollection.InsertOne(context.TODO(), notification) // Insert the notification
	if err != nil {                                                          // Check if there is an error
		// Log the error and return it
		log.Printf("Failed to insert session notification: %v\n", err) // Log the error message
		return err                                                     // Return the error
	}

	// Log success for debugging
	log.Println("Session notification sent successfully") // Log the success message
	return nil                                            // Return nil
}

// User Notification: Sends notifications to users about invitations, changes, and updates.
// SendUserNotification handles sending notifications to users about invitations, changes, and updates.
func SendUserNotification(c *gin.Context) { // Send a notification to a user
	var notification models.Notification // Define a notification variable

	// Bind the JSON payload to the notification struct
	if err := c.BindJSON(&notification); err != nil { // Bind the JSON to the notification struct
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return an error response
		return                                                         // Return from the function
	}

	// Validate if the UserID exists in the user collection
	var user models.User                                                                            // Define a user variable
	err := userCollection.FindOne(context.TODO(), bson.M{"_id": notification.UserID}).Decode(&user) // Find the user by ID
	if err != nil {                                                                                 // Check if there is an error
		// If user does not exist, return an error response
		if err == mongo.ErrNoDocuments { // Check if the user was not found
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"}) // Return an error response
		} else { // Handle other errors
			// Handle other errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		}
		return // Return from the function
	}

	// Set additional fields for the notification
	notification.ID = primitive.NewObjectID() // Generate a new ObjectID for the notification
	notification.Type = "User"                // Set the type of notification
	notification.CreatedAt = time.Now()       // Set the created_at timestamp
	notification.UpdatedAt = time.Now()       // Set the updated_at timestamp

	// Insert the notification into the database
	_, err = notificationCollection.InsertOne(context.TODO(), notification) // Insert the notification
	if err != nil {                                                         // Check if there is an error
		// Handle insertion errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "User notification sent successfully"}) // Return the created notification
}

// Get Notifications: Allows users to view their notifications.
func GetNotifications(c *gin.Context) { // Get notifications for a user
	userID := c.Param("userId") // Get user ID from the URL

	objectUserID, err := primitive.ObjectIDFromHex(userID) // Convert user ID to ObjectID
	if err != nil {                                        // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return an error response
		return                                                           // Return from the function
	}

	var notifications []models.Notification // Define a notifications variable

	cursor, err := notificationCollection.Find(context.TODO(), bson.M{"userId": objectUserID}) // Find notifications by user ID
	if err != nil {                                                                            // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor

	for cursor.Next(context.TODO()) { // Iterate over the cursor
		var notification models.Notification                 // Define a notification variable
		if err := cursor.Decode(&notification); err != nil { // Decode the notification
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
			return                                                              // Return from the function
		}
		notifications = append(notifications, notification) // Append the notification to the notifications slice
	}

	c.JSON(http.StatusOK, notifications) // Return the notifications
}

// Delete Notification: Allows users to delete a notification.
func DeleteNotification(c *gin.Context) { // Delete a notification
	notificationID := c.Param("notificationId") // Get the notification ID from the URL

	objectNotificationID, err := primitive.ObjectIDFromHex(notificationID) // Convert ID to ObjectID
	if err != nil {                                                        // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"}) // Return an error response
		return                                                                   // Return from the function
	}

	_, err = notificationCollection.DeleteOne(context.TODO(), bson.M{"_id": objectNotificationID}) // Delete the notification
	if err != nil {                                                                                // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"}) // Return a success response
}
