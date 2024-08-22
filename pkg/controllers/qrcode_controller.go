package controllers

import (
	"context"
	"fmt"
	"net/http"
	"training_session/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ( // Define a variable for the session collection in QR code controller
	qrSessionCollection *mongo.Collection // Use a different variable for the session collection in QR code controller
)

func InitializeQRCodeController(database *mongo.Database) { // Initialize the QR code controller
	qrSessionCollection = database.Collection("sessions") // Set the session collection for QR code controller
}

// GenerateQRCode generates a QR code for session verification
func GenerateQRCode(c *gin.Context) { // Generate a QR code for session verification
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	// Convert sessionID to ObjectID
	_, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"}) // Return an error response
		return                                                              // Return from the function to stop execution
	}

	// Generate QR code content
	qrContent := fmt.Sprintf("session:%s", sessionID) // Generate the QR code content

	// Create QR code
	code, err := qrcode.Encode(qrContent, qrcode.Medium, 256) // Generate the QR code
	if err != nil {                                           // Check if there is an error generating the QR code
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"}) // Return an error response
		return                                                                               // Return from the function to stop execution
	}

	// Return QR code image
	c.Header("Content-Type", "image/png") // Set the content type to image/png
	c.Writer.WriteHeader(http.StatusOK)   // Set the status code to 200
	c.Writer.Write(code)                  // Write the QR code image to the response
} // GenerateQRCode generates a QR code for session verification

// ValidateQRCode validates the QR code to ensure session integrity
func ValidateQRCode(c *gin.Context) {
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	// Convert sessionID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                       // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"}) // Return an error response
		return                                                              // Return from the function to stop execution
	}

	// Find the session in the database
	var session models.Session                                                                // Define a session variable
	err = sessionCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&session) // Find the session by ID
	if err != nil {                                                                           // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the session was not found
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"}) // Return a not found response
		} else { // If there is another error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		}
		return // Return from the function to stop execution
	}

	// Check if the session is valid (you can add more validation logic here)
	if session.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is not active"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "QR code is valid"}) // Return a success response
}
