package controllers

import (
	"context"
	"net/http"
	"time"
	"training_session/pkg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	pitchBookingCollection *mongo.Collection // Define a pitchBookingCollection variable
)

// InitializePitchBookingController initializes the pitch booking controller
func InitializePitchBookingController(database *mongo.Database) {
	pitchBookingCollection = database.Collection("pitch_bookings") // Set the pitch booking collection
}

// GetPitchBookings retrieves all pitch bookings
func GetPitchBookings(c *gin.Context) {
	var pitchBookings []models.Pitch                                       // Define a pitchBookings variable
	cursor, err := pitchBookingCollection.Find(context.TODO(), bson.D{{}}) // Find all pitch bookings
	if err != nil {                                                        // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor
	for cursor.Next(context.TODO()) {  // Iterate over the cursor
		var pitchBooking models.Pitch                        // Define a pitchBooking variable
		if err := cursor.Decode(&pitchBooking); err != nil { // Decode the pitch booking
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
			return                                                              // Return from the function
		}
		pitchBookings = append(pitchBookings, pitchBooking) // Append the pitch booking to the pitchBookings slice
	}
	c.JSON(http.StatusOK, pitchBookings) // Return a success response
}

// GetPitchBookingByID retrieves a pitch booking by ID
func GetPitchBookingByID(c *gin.Context) {
	pitchBookingID := c.Param("pitchId") // Get the pitch booking ID from the URL

	objectID, err := primitive.ObjectIDFromHex(pitchBookingID) // Convert ID to ObjectID
	if err != nil {                                            // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pitch booking ID"}) // Return a bad request response
		return                                                                    // Return from the function
	}

	var pitchBooking models.Pitch                                                                       // Define a pitchBooking variable
	err = pitchBookingCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&pitchBooking) // Find and decode the pitch booking
	if err != nil {                                                                                     // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if no documents were found
			c.JSON(http.StatusNotFound, gin.H{"error": "Pitch booking not found"}) // Return a not found response
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		}
		return
	}

	c.JSON(http.StatusOK, pitchBooking) // Return the pitch booking
}

// CreatePitchBooking creates a new pitch booking
func BookPitch(c *gin.Context) {
	var pitchBooking models.Pitch                     // Define a pitchBooking variable
	if err := c.BindJSON(&pitchBooking); err != nil { // Bind the JSON to the pitch booking struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return a bad request response
		return                                                     // Return from the function
	}

	pitchBooking.ID = primitive.NewObjectID() // Generate a new ObjectID for the pitch booking
	pitchBooking.CreatedAt = time.Now()       // Set the created_at timestamp
	pitchBooking.UpdatedAt = time.Now()       // Set the updated_at timestamp

	_, err := pitchBookingCollection.InsertOne(context.TODO(), pitchBooking) // Insert the pitch booking
	if err != nil {                                                          // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pitch booking"}) // Return an error response
		return                                                                                   // Return from the function
	}

	c.JSON(http.StatusCreated, pitchBooking) // Return the created pitch booking
}

// UpdatePitchBooking updates an existing pitch booking
func UpdatePitchBooking(c *gin.Context) {
	pitchBookingID := c.Param("pitchID") // Get the pitch booking ID from the URL
	var updatedPitchBooking models.Pitch // Define an updated pitch booking variable

	if err := c.BindJSON(&updatedPitchBooking); err != nil { // Bind the JSON to the updated pitch booking struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return a bad request response
		return                                                     // Return from the function
	}

	objectID, err := primitive.ObjectIDFromHex(pitchBookingID) // Convert ID to ObjectID
	if err != nil {                                            // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pitch booking ID"}) // Return a bad request response
		return                                                                    // Return from the function
	}

	updatedPitchBooking.UpdatedAt = time.Now() // Set the updated_at timestamp

	_, err = pitchBookingCollection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, updatedPitchBooking) // Update the pitch booking
	if err != nil {                                                                                          // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pitch booking"}) // Return an error response
		return                                                                                   // Return from the function
	}

	c.JSON(http.StatusOK, updatedPitchBooking) // Return the updated pitch booking
}

// DeletePitchBooking deletes an existing pitch booking
func DeletePitchBooking(c *gin.Context) {
	pitchBookingID := c.Param("pitchId") // Get the pitch booking ID from the URL

	objectID, err := primitive.ObjectIDFromHex(pitchBookingID) // Convert ID to ObjectID
	if err != nil {                                            // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pitch booking ID"}) // Return a bad request response
		return                                                                    // Return from the function
	}

	_, err = pitchBookingCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID}) // Delete the pitch booking
	if err != nil {                                                                    // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pitch booking"}) // Return an error response
		return                                                                                   // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pitch booking deleted successfully"}) // Return a success response
}

// GetPitchBookingsByPitchID retrieves all pitch bookings by pitch ID
func GetPitchBookingsByPitchID(c *gin.Context) {
	pitchID := c.Param("pitchId") // Get the pitch ID from the URL

	objectPitchID, err := primitive.ObjectIDFromHex(pitchID) // Convert ID to ObjectID
	if err != nil {                                          // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pitch ID"}) // Return a bad request response
		return                                                            // Return from the function
	}

	var pitchBookings []models.Pitch                                                             // Define a pitchBookings variable
	cursor, err := pitchBookingCollection.Find(context.TODO(), bson.M{"pitchId": objectPitchID}) // Find pitch bookings by pitch ID
	if err != nil {                                                                              // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor
	for cursor.Next(context.TODO()) {  // Iterate over the cursor
		var pitchBooking models.Pitch                        // Define a pitchBooking variable
		if err := cursor.Decode(&pitchBooking); err != nil { // Decode the pitch booking
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
			return                                                              // Return from the function
		}
		pitchBookings = append(pitchBookings, pitchBooking) // Append the pitch booking to the pitchBookings slice
	}

	c.JSON(http.StatusOK, pitchBookings) // Return a success response
}

// GetPitchBookingsByUserID retrieves all pitch bookings by user ID
func GetPitchBookingsByUserID(c *gin.Context) {
	userID := c.Param("userId") // Get the user ID from the URL

	objectUserID, err := primitive.ObjectIDFromHex(userID) // Convert ID to ObjectID
	if err != nil {                                        // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return a bad request response
		return                                                           // Return from the function
	}

	var pitchBookings []models.Pitch                                                           // Define a pitchBookings variable
	cursor, err := pitchBookingCollection.Find(context.TODO(), bson.M{"userId": objectUserID}) // Find pitch bookings by user ID
	if err != nil {                                                                            // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor
	for cursor.Next(context.TODO()) {  // Iterate over the cursor
		var pitchBooking models.Pitch                        // Define a pitchBooking variable
		if err := cursor.Decode(&pitchBooking); err != nil { // Decode the pitch booking
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
			return                                                              // Return from the function
		}
		pitchBookings = append(pitchBookings, pitchBooking) // Append the pitch booking to the pitchBookings slice
	}

	c.JSON(http.StatusOK, pitchBookings) // Return a success response
}

// GetPitchBookingsByDate retrieves all pitch bookings by date
func GetPitchBookingsByDate(c *gin.Context) {
	date := c.Param("date") // Get the date from the URL

	var pitchBookings []models.Pitch                                                 // Define a pitchBookings variable
	cursor, err := pitchBookingCollection.Find(context.TODO(), bson.M{"date": date}) // Find pitch bookings by date
	if err != nil {                                                                  // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor
	for cursor.Next(context.TODO()) {  // Iterate over the cursor
		var pitchBooking models.Pitch                        // Define a pitchBooking variable
		if err := cursor.Decode(&pitchBooking); err != nil { // Decode the pitch booking
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
			return                                                              // Return from the function
		}
		pitchBookings = append(pitchBookings, pitchBooking) // Append the pitch booking to the pitchBookings slice
	}

	c.JSON(http.StatusOK, pitchBookings) // Return a success response
}
