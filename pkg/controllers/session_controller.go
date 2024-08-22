package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"training_session/pkg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var sessionCollection *mongo.Collection // Define a sessionCollection variable

func InitializeSession(db *mongo.Database) { // Initialize the controllers
	sessionCollection = db.Collection("sessions") // Set the session collection
}

func GetSessions(c *gin.Context) { // Get all sessions
	var sessions []models.Session // Define a sessions variable

	cursor, err := sessionCollection.Find(context.TODO(), bson.D{{}}) // Find all sessions
	if err != nil {                                                   // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
			"error": err.Error(), // Return the error message
		})
		return // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor

	for cursor.Next(context.TODO()) { // Iterate over the cursor
		var session models.Session                      // Define a session variable
		if err := cursor.Decode(&session); err != nil { // Decode the session
			c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
				"error": err.Error(), // Return the error message
			})
			return // Return from the function
		}
		sessions = append(sessions, session) // Append the session to the sessions slice
	}

	c.JSON(http.StatusOK, sessions) // Return a success response
}

func GetSessionByID(c *gin.Context) { // Get a session by ID
	sessionID := c.Param("id") // Get the session ID from the URL

	// Convert the sessionID from string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                       // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid session ID format", // Return an error message
		})
		return // Return from the function
	}

	var session models.Session                                                   // Define a session variable
	result := sessionCollection.FindOne(context.TODO(), bson.M{"_id": objectID}) // Find the session by ID
	if err := result.Err(); err != nil {                                         // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the session was not found
			c.JSON(http.StatusNotFound, gin.H{ // Return a not found response
				"error": "Session not found", // Return an error message
			})
		} else { // Check if there is another error
			c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
				"error": err.Error(), // Return the error message
			})
		}
		return
	}

	if err := result.Decode(&session); err != nil { // Decode the session
		c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
			"error": err.Error(), // Return the error message
		})
		return // Return from the function
	}

	c.JSON(http.StatusOK, session) // Return the session
}

func CreateSession(c *gin.Context) { // Create a session
	if sessionCollection == nil { // Check if the session collection is not initialized
		c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
			"error": "MongoDB session collection is not initialized", // Return an error message
		})
		return // Return from the function
	}

	var session models.Session                         // Define a session variable
	if err := c.ShouldBindJSON(&session); err != nil { // Bind the JSON data to the session variable
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": err.Error(), // Return the error message
		})
		return
	}

	session.CreatedAt = time.Now() // Set the created time
	session.UpdatedAt = time.Now() // Set the updated time

	_, err := sessionCollection.InsertOne(context.TODO(), session) // Insert the session
	if err != nil {                                                // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
			"error": fmt.Sprintf("Failed to create session: %v", err), // Return the error message
		})
		return
	}

	c.JSON(http.StatusCreated, session) // Return the created session
}

func UpdateSession(c *gin.Context) { // Update a session
	sessionID := c.Param("id") // Get the session ID from the URL

	var session models.Session // Define a session variable

	if err := c.ShouldBindJSON(&session); err != nil { // Bind the JSON data to the session variable
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": err.Error(), // Return the error message
		})
		return
	}

	// Convert the sessionID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                       // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid session ID format", // Return an error message
		})
		return
	}

	session.ID = objectID          // Set the session ID to the converted ObjectID
	session.UpdatedAt = time.Now() // Set the updated time

	// Update the session
	result, err := sessionCollection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, session) // Replace the session
	if err != nil {                                                                               // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
			"error": err.Error(), // Return the error message
		})
		return // Return from the function
	}

	// Check if any document was modified
	if result.ModifiedCount == 0 { // Check if the document was not modified
		c.JSON(http.StatusNotFound, gin.H{ // Return a not found response
			"error": "No document found with the given ID", // Return an error message
		})
		return // Return from the function
	}

	c.JSON(http.StatusOK, session) // Return the updated session
}

func DeleteSession(c *gin.Context) { // Delete a session
	sessionID := c.Param("id") // Get the session ID from the URL

	// Convert sessionID from string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid session ID format", // Return an error message
		})
		return // Return from the function
	}

	result, err := sessionCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID}) // Delete the session
	if err != nil {                                                                     // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{ // Return an error response
			"error": err.Error(), // Return the error message
		})
		return // Return from the function
	}

	if result.DeletedCount == 0 { // Check if the session was not found
		c.JSON(http.StatusNotFound, gin.H{ // Return a not found response
			"error": "Session not found", // Return an error message
		})
		return // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{ // Return a success response
		"message": "Session deleted", // Return a message
	})
}
