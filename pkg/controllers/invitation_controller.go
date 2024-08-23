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

var invitationCollection *mongo.Collection // Define an invitationCollection variable

func InitializeInvitation(database *mongo.Database) { // Initialize the controllers
	invitationCollection = database.Collection("invitations") // Set the invitation collection
}

func SendInvitation(c *gin.Context) { // Send an invitation for a private training session
	var invitation models.Invitation                // Define an invitation variable
	if err := c.BindJSON(&invitation); err != nil { // Bind the JSON to the invitation struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Return a bad request response
		return                                                     // Return from the function
	}

	invitation.ID = primitive.NewObjectID() // Generate a new ObjectID for the invitation
	invitation.CreatedAt = time.Now()       // Set the created_at timestamp
	invitation.UpdatedAt = time.Now()       // Set the updated_at timestamp

	_, err := invitationCollection.InsertOne(context.TODO(), invitation) // Insert the invitation
	if err != nil {                                                      // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusCreated, invitation) // Return the created invitation
}

func AcceptInvitation(c *gin.Context) { // Handle user acceptance of session invitations
	invitationID := c.Param("invitationId") // Get the invitation ID from the URL

	// Convert invitationID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(invitationID) // Convert the invitation ID to an ObjectID
	if err != nil {                                          // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"}) // Return a bad request response
		return                                                                 // Return from the function
	}

	// Update the invitation document to mark as accepted
	filter := bson.M{"_id": objectID}                                             // Define the filter
	update := bson.M{"$set": bson.M{"status": "accepted"}}                        // Define the update
	result, err := invitationCollection.UpdateOne(context.TODO(), filter, update) // Update the invitation
	if err != nil {                                                               // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	if result.MatchedCount == 0 { // Check if the invitation was not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"}) // Return a not found response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted"}) // Return a success response
}

func DeclineInvitation(c *gin.Context) { // Manage user decline of session invitations
	invitationID := c.Param("invitationId") // Get the invitation ID from the URL

	// Convert invitationID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(invitationID) // Convert the invitation ID to an ObjectID
	if err != nil {                                          // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"}) // Return a bad request response
		return                                                                 // Return from the function
	}

	// Update the invitation document to mark as declined
	filter := bson.M{"_id": objectID}                                             // Define the filter
	update := bson.M{"$set": bson.M{"status": "declined"}}                        // Define the update
	result, err := invitationCollection.UpdateOne(context.TODO(), filter, update) // Update the invitation
	if err != nil {                                                               // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	if result.MatchedCount == 0 { // Check if the invitation was not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"}) // Return a not found response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation declined"}) // Return a success response
}

func GetInvitations(c *gin.Context) { // Get all invitations
	var invitations []models.Invitation                                  // Define a slice to hold the invitations
	cursor, err := invitationCollection.Find(context.TODO(), bson.D{{}}) // Find all invitations
	if err != nil {                                                      // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return
	}
	defer cursor.Close(context.TODO()) // Close the cursor
	for cursor.Next(context.TODO()) {  // Iterate over the cursor
		var invitation models.Invitation                   // Define an invitation variable
		if err := cursor.Decode(&invitation); err != nil { // Decode the invitation
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
			return
		}
		invitations = append(invitations, invitation) // Append the invitation to the slice
	}
	c.JSON(http.StatusOK, invitations) // Return a success response
}

func GetInvitationByID(c *gin.Context) { // Get an invitation by ID
	invitationID := c.Param("invitationId") // Get the invitation ID from the URL

	objectID, err := primitive.ObjectIDFromHex(invitationID) // Convert ID to ObjectID
	if err != nil {                                          // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"}) // Return a bad request response
		return
	}

	var invitation models.Invitation                                                                // Define an invitation variable
	err = invitationCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&invitation) // Find and decode the invitation
	if err != nil {                                                                                 // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if no documents were found
			c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"}) // Return a not found response
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		}
		return
	}

	c.JSON(http.StatusOK, invitation) // Return the invitation
}

func DeleteInvitation(c *gin.Context) { // Delete an invitation
	invitationID := c.Param("invitationId") // Get the invitation ID from the URL

	objectID, err := primitive.ObjectIDFromHex(invitationID) // Convert ID to ObjectID
	if err != nil {                                          // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invitation ID"}) // Return a bad request response
		return
	}

	result, err := invitationCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID}) // Delete the invitation
	if err != nil {                                                                        // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return
	}

	if result.DeletedCount == 0 { // Check if no document was deleted
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found"}) // Return a not found response
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation deleted successfully"}) // Return a success response
}
