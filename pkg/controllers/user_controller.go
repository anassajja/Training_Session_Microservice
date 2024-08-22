package controllers

import (
	"context"
	"net/http"
	"training_session/pkg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userCollection *mongo.Collection // Define a userCollection variable
)

func InitializeUser(database *mongo.Database) { // Initialize the controllers
	userCollection = database.Collection("users")       // Set the user collection
	sessionCollection = database.Collection("sessions") // Set the session collection
}

func GetUsers(c *gin.Context) { // Get all users
	var users []models.User                                        // Define a users variable
	cursor, err := userCollection.Find(context.TODO(), bson.D{{}}) // Find all users
	if err != nil {                                                // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor
	for cursor.Next(context.TODO()) {  // Iterate over the cursor
		var user models.User                         // Define a user variable
		if err := cursor.Decode(&user); err != nil { // Decode the user
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
			return                                                              // Return from the function
		}
		users = append(users, user) // Append the user to the users slice
	}
	c.JSON(http.StatusOK, users) // Return a success response
}

func RegisterUser(c *gin.Context) { // Create a user
	var user models.User                      // Define a user variable
	if err := c.BindJSON(&user); err != nil { // Bind the JSON to the user struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return a bad request response
		return                                                     // Return from the function
	}
	_, err := userCollection.InsertOne(context.TODO(), user) // Insert the user
	if err != nil {                                          // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	c.JSON(http.StatusCreated, user) // Return the created user
}

func GetUserByID(c *gin.Context) { // Get a user by ID
	userID := c.Param("id") // Get the user ID from the URL

	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID) // Convert the user ID to an ObjectID
	if err != nil {                                    // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return a bad request response
		return                                                           // Return from the function
	}

	// Find the user document
	var user models.User                                                                // Define a user variable
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user) // Find the user by ID
	if err != nil {                                                                     // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the user was not found
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return a not found response
			return                                                        // Return from the function
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusOK, user) // Return the user
}

func UpdateUser(c *gin.Context) { // Update a user
	userID := c.Param("id") // Get the user ID from the URL
	var user models.User    // Define a user variable

	// Bind JSON to user struct
	if err := c.BindJSON(&user); err != nil { // Bind the JSON to the user struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return a bad request response
		return                                                     // Return from the function
	}

	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID) // Convert the user ID to an ObjectID
	if err != nil {                                    // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return a bad request response
		return                                                           // Return from the function
	}

	// Update the user document
	filter := bson.M{"_id": objectID}                                       // Define the filter
	update := bson.M{"$set": user}                                          // Define the update
	result, err := userCollection.UpdateOne(context.TODO(), filter, update) // Update the user
	if err != nil {                                                         // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	if result.MatchedCount == 0 { // Check if the user was not found
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return a not found response
		return                                                        // Return from the function
	}

	c.JSON(http.StatusOK, user) // Return the updated user
}

func DeleteUser(c *gin.Context) { // Delete a user
	userID := c.Param("id") // Get the user ID from the URL

	// Convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID) // Convert the user ID to an ObjectID
	if err != nil {                                    // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return a bad request response
		return                                                           // Return from the function
	}

	// Delete the user document
	filter := bson.M{"_id": objectID}                               // Define the filter
	result, err := userCollection.DeleteOne(context.TODO(), filter) // Delete the user
	if err != nil {                                                 // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	if result.DeletedCount == 0 { // Check if the user was not found
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return a not found response
		return                                                        // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"}) // Return a success response
}

func EnrollInSession(c *gin.Context) { // Enroll user in a session
	userID := c.Param("userId")       // Get the user ID from the URL
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	// Check if user exists
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user) // Find the user by ID
	if err != nil {                                                                    // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the user was not found
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // Return a not found response
			return                                                        // Return from the function
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	// Check if session exists
	var session models.Session                                                                 // Define a session variable
	err = sessionCollection.FindOne(context.TODO(), bson.M{"_id": sessionID}).Decode(&session) // Find the session by ID
	if err != nil {                                                                            // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the session was not found
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"}) // Return a not found response
			return                                                           // Return from the function
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	// Check if user is already enrolled
	for _, participant := range session.Participants { // Iterate over the participants array
		if participant == userID { // Check if the user is already enrolled
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already enrolled in the session"}) // Return a bad request response
			return                                                                                // Return from the function
		}
	}

	// Enroll user in the session
	_, err = sessionCollection.UpdateOne( // Update the session document to add the user to the participants array field
		context.TODO(),           // Context for the operation to execute in the database server environment (required by the driver)
		bson.M{"_id": sessionID}, // Filter to find the session by ID
		bson.M{"$push": bson.M{"participants": userID}}, // Update operation to push the user ID to the participants array
	)
	if err != nil { // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User enrolled in session successfully"}) // Return a success response
}

func ProcessRefund(c *gin.Context) { // Process a refund for a user
	userID := c.Param("userId")       // Get the user ID from the URL
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	// Convert userID and sessionID to ObjectID
	objectUserID, err := primitive.ObjectIDFromHex(userID) // Convert the user ID to an ObjectID
	if err != nil {                                        // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return a bad request response
		return                                                           // Return from the function
	}

	objectSessionID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                              // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"}) // Return a bad request response
		return
	}

	// Implement your logic to process the refund here
	// For example, update the payment status in your payment service or call an external payment API

	// Placeholder for refund processing logic
	// You should integrate with a payment gateway or service to handle actual refunds
	// Example: PaymentService.Refund(userID, sessionID)

	c.JSON(http.StatusOK, gin.H{"message": "Refund processed successfully"}) // Return a success response
}

func CancelEnrollment(c *gin.Context) { // Cancel user enrollment in a session
	userID := c.Param("userId")       // Get the user ID from the URL
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	// Convert userID and sessionID to ObjectID
	objectUserID, err := primitive.ObjectIDFromHex(userID) // Convert the user ID to an ObjectID
	if err != nil {                                        // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return a bad request response
		return                                                           // Return from the function
	}

	objectSessionID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                              // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"}) // Return a bad request response
		return                                                              // Return from the function
	}

	// Remove user from the session
	filter := bson.M{"_id": objectSessionID}                                   // Define the filter
	update := bson.M{"$pull": bson.M{"participants": objectUserID}}            // Define the update operation
	result, err := sessionCollection.UpdateOne(context.TODO(), filter, update) // Update the session
	if err != nil {                                                            // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	if result.MatchedCount == 0 { // Check if the session was not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Session or user not found"}) // Return a not found response
		return                                                                   // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment canceled successfully"}) // Return a success response
}

func CancelSession(c *gin.Context) { // Cancel a session
	sessionID := c.Param("id") // Get the session ID from the URL

	// Convert sessionID to ObjectID to find the session
	objectSessionID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                              // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"}) // Return a bad request response
		return                                                              // Return from the function
	}

	// Find the session to get participants
	var session models.Session                                                                       // Define a session variable
	err = sessionCollection.FindOne(context.TODO(), bson.M{"_id": objectSessionID}).Decode(&session) // Find the session by ID
	if err != nil {                                                                                  // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the session was not found
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"}) // Return a not found response
		} else { // Check if there is another error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		}
		return // Return from the function
	}

	// Notify participants about the cancellation
	// Example: Send email notifications or any other form of notification
	for _, participantID := range session.Participants { // Iterate over the participants
		// Here, implement the logic to notify each participant
		// For example, you could send an email to each participant
		// or push a notification using a service
		_ = participantID // Replace this with actual notification logic
	}

	// Delete the session
	_, err = sessionCollection.DeleteOne(context.TODO(), bson.M{"_id": objectSessionID}) // Delete the session
	if err != nil {                                                                      // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session canceled successfully"}) // Return a success response
}
