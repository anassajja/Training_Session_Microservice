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

var feedbackCollection *mongo.Collection

func InitializeFeedbackController(database *mongo.Database) { // Initialize the feedback controller
	feedbackCollection = database.Collection("feedbacks") // Set the feedback collection
}

// SubmitFeedback: Manages the submission of feedback for sessions and coaches
func SubmitFeedback(c *gin.Context) { // Submit feedback for sessions and coaches
	var feedback models.Feedback                  // Define a feedback variable
	if err := c.BindJSON(&feedback); err != nil { // Bind the JSON to the feedback struct
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return a bad request response
		return                                                         // Return from the function
	}

	feedback.ID = primitive.NewObjectID() // Generate a new ObjectID for the feedback
	feedback.CreatedAt = time.Now()       // Set the created_at timestamp
	feedback.UpdatedAt = time.Now()       // Set the updated_at timestamp

	_, err := feedbackCollection.InsertOne(context.TODO(), feedback) // Insert the feedback
	if err != nil {                                                  // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit feedback"}) // Return an error response
		return                                                                              // Return from the function
	}

	c.JSON(http.StatusCreated, feedback) // Return the created feedback
}

// ViewFeedback: Allows users to view feedback they have submitted
func ViewFeedback(c *gin.Context) { // View feedback submitted by a user
	userID := c.Param("userId") // Get user ID from the URL

	objectUserID, err := primitive.ObjectIDFromHex(userID) // Convert user ID to ObjectID
	if err != nil {                                        // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"}) // Return an error response
		return                                                           // Return from the function
	}

	// Find feedbacks submitted by this user
	var feedbacks []models.Feedback                                                         // Define a feedbacks variable
	cursor, err := feedbackCollection.Find(context.TODO(), bson.M{"user_id": objectUserID}) // Find feedbacks by user ID
	if err != nil {                                                                         // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve feedback"}) // Return an error response
		return                                                                                // Return from the function
	}
	defer cursor.Close(context.TODO()) // Close the cursor

	for cursor.Next(context.TODO()) { // Iterate over the cursor
		var feedback models.Feedback                     // Define a feedback variable
		if err := cursor.Decode(&feedback); err != nil { // Decode the feedback
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding feedback"}) // Return an error response
			return                                                                            // Return from the function
		}
		feedbacks = append(feedbacks, feedback) // Append the feedback to the feedbacks slice
	}

	c.JSON(http.StatusOK, feedbacks) // Return the feedbacks
}

// EditFeedback: Handles editing of previously submitted feedback
func EditFeedback(c *gin.Context) { // Edit previously submitted feedback
	feedbackID := c.Param("id")         // Get feedback ID from the URL
	var updatedFeedback models.Feedback // Define an updated feedback variable

	if err := c.BindJSON(&updatedFeedback); err != nil { // Bind the JSON to the updated feedback struct
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"}) // Return a bad request response
		return                                                         // Return from the function
	}

	objectFeedbackID, err := primitive.ObjectIDFromHex(feedbackID) // Convert feedback ID to ObjectID
	if err != nil {                                                // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"}) // Return an error response
		return                                                               // Return from the function
	}

	updatedFeedback.UpdatedAt = time.Now() // Update the updated_at timestamp

	// Find and update feedback
	filter := bson.M{"_id": objectFeedbackID} // Define the filter
	update := bson.M{"$set": bson.M{          // Define the update
		"content":    updatedFeedback.Content,   // Define the update
		"rating":     updatedFeedback.Rating,    // Update the content and rating
		"updated_at": updatedFeedback.UpdatedAt, // Update the updated_at timestamp
	}}

	result, err := feedbackCollection.UpdateOne(context.TODO(), filter, update) // Update the feedback
	if err != nil {                                                             // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit feedback"}) // Return an error response
		return                                                                            // Return from the function
	}

	if result.MatchedCount == 0 { // Check if the feedback was not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"}) // Return a not found response
		return                                                            // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feedback updated successfully"}) // Return a success response
}

// DeleteFeedback: Manages deletion of feedback if necessary
func DeleteFeedback(c *gin.Context) { // Delete feedback
	feedbackID := c.Param("id") // Get feedback ID from the URL

	objectFeedbackID, err := primitive.ObjectIDFromHex(feedbackID) // Convert feedback ID to ObjectID
	if err != nil {                                                // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"}) // Return an error response
		return                                                               // Return from the function
	}

	// Find and delete feedback
	filter := bson.M{"_id": objectFeedbackID}                           // Define the filter
	result, err := feedbackCollection.DeleteOne(context.TODO(), filter) // Delete the feedback
	if err != nil {                                                     // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete feedback"}) // Return an error response
		return                                                                              // Return from the function
	}

	if result.DeletedCount == 0 { // Check if the feedback was not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"}) // Return a not found response
		return                                                            // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feedback deleted successfully"}) // Return a success response
}
