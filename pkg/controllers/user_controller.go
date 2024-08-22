package controllers

import (
	"context"
	"net/http"
	"training_session/db"
	"training_session/pkg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo" // Import the MongoDB driver
)

var userCollection *mongo.Collection // Define a userCollection variable

// Initialize the user collection only after MongoDB client is connected
func SetupUserCollection() { // Set up the user collection
	if db.Client == nil { // Check if the client is not initialized
		panic("MongoDB client is not initialized") // Panic if the client is not initialized
	}
	userCollection = db.Client.Database("test").Collection("users") // Set the userCollection variable
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

func CreateUser(c *gin.Context) { // Create a user
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
