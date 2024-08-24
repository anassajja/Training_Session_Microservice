package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"training_session/config"
	"training_session/pkg/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	sessionCollection *mongo.Collection // Define a sessionCollection variable
)

func InitializeSession(db *mongo.Database) { // Initialize the controllers
	sessionCollection = db.Collection("sessions") // Set the session collection
}

func init() {
	cfg = config.LoadConfig()
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

func GetActiveSessions(c *gin.Context) { // Get all active sessions
	var sessions []models.Session // Define a sessions variable

	cursor, err := sessionCollection.Find(context.TODO(), bson.M{"status": "active"}) // Find all active sessions
	if err != nil {                                                                   // Check if there is an error
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
	sessionID := c.Param("sessionId") // Get the session ID from the URL

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

func GetSessionsByUserID(c *gin.Context) { // Get all sessions created by a user
	userID := c.Param("userId") // Get the user ID from the URL

	// Convert the userID from string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID) // Convert the user ID to an ObjectID
	if err != nil {                                    // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid user ID format", // Return an error message
		})
		return // Return from the function
	}

	var sessions []models.Session // Define a sessions variable

	cursor, err := sessionCollection.Find(context.TODO(), bson.M{"user_id": objectID}) // Find all sessions by user ID
	if err != nil {                                                                    // Check if there is an error
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

func CreateSession(c *gin.Context) {
	if sessionCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "MongoDB session collection is not initialized",
		})
		return
	}
	log.Printf("JWT Secret: %s", cfg.JwtSecretKey)

	tokenString, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtSecretKey), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userIDHex, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User // Define a user variable to store the user details from the database collection
	err = userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user) // Find the user by ID and decode the result
	if err != nil { 																 // Check if there is an error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"}) // Return an unauthorized response with an error message
		return 																	 // Return from the function
	}

	if user.Role != "coach" && user.Role != "business owner" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You do not have the required permissions to create a session",
		})
		return
	}

	var session models.Session
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	session.ID = primitive.NewObjectID()
	session.Status = "active"
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	_, err = sessionCollection.InsertOne(context.TODO(), session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create session: %v", err),
		})
		return
	}

	notification := models.Notification{
		ID:        primitive.NewObjectID(),
		UserID:    user.ID,
		Type:      "Session Created",
		Message:   fmt.Sprintf("The session '%s' has been created.", session.Title),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = SendSessionNotification(notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// UpdateSession: Updates a session and sends a notification to the creator
func UpdateSession(c *gin.Context) { // Update a session
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	var session models.Session // Define a session variable
	var user models.User       // Define a user variable

	if err := c.ShouldBindJSON(&session); err != nil { // Bind the JSON data to the session variable
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": err.Error(), // Return the error message
		}) // Return a bad request response
		return // Return from the function
	}

	// Convert the sessionID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                       // Check if there is an error
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid session ID format", // Return an error message
		})
		return // Return from the function
	}

	session.ID = objectID          // Set the session ID to the converted ObjectID
	session.UpdatedAt = time.Now() // Set the updated time

	// Update the session
	result, err := sessionCollection.ReplaceOne(context.TODO(), bson.M{"_id": objectID}, session) // Replace the session document with the updated session
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
		return
	}

	// Prepare notification data
	notification := models.Notification{ // Define a notification variable with the required fields
		ID:        primitive.NewObjectID(),                                          // Generate a new ObjectID for the notification
		UserID:    user.ID,                                                          // Assuming session.CreatorID holds the ID of the user to notify
		Type:      "Session Updated",                                                // Set the notification type
		Message:   fmt.Sprintf("The session '%s' has been updated.", session.Title), // Message for the notification
		CreatedAt: time.Now(),                                                       // Set the created_at timestamp
		UpdatedAt: time.Now(),                                                       // Set the updated_at timestamp
	}

	// Send session notification
	err = SendSessionNotification(notification) // Send the session notification
	if err != nil {                             // Check if there is an error sending the notification
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"}) // Return an error response
		return                                                                                // Return from the function
	}

	c.JSON(http.StatusOK, session) // Return the updated session
}

// DeleteSession: Deletes a session from the database
func DeleteSession(sessionID primitive.ObjectID) error { // Delete a session
	_, err := sessionCollection.DeleteOne(context.TODO(), bson.M{"_id": sessionID}) // Delete the session
	return err                                                                      // Return the error
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
	filter := bson.M{"_id": objectSessionID}                                   // Define the filter to find the session by ID
	update := bson.M{"$pull": bson.M{"participants": objectUserID}}            // Define the update operation to remove the user from the participants array
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

// CancelSession: Cancels a session and sends a notification to the creator
func CancelSession(c *gin.Context) { // Cancel a session
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	// Convert sessionID to ObjectID
	objectSessionID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                              // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"}) // Return a bad request response
		return                                                              // Return from the function
	} // Check if there is an error converting the ID

	// Find the session to get details
	var session models.Session
	var user models.User                                                                             // Define a session variable
	err = sessionCollection.FindOne(context.TODO(), bson.M{"_id": objectSessionID}).Decode(&session) // Find the session by ID
	if err != nil {                                                                                  // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the session was not found
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"}) // Return a not found response
		} else { // Check if there is another error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		} // Check if there is another error
		return // Return from the function
	}

	// Prepare notification data
	notification := models.Notification{ // Define a notification variable with the required fields
		ID:        primitive.NewObjectID(),                                           // Generate a new ObjectID for the notification
		UserID:    user.ID,                                                           // Assuming this holds the ID of the user to notify
		Type:      "Session Cancellation",                                            // Set the type of notification
		Message:   fmt.Sprintf("The session '%s' has been canceled.", session.Title), // Set the message
		CreatedAt: time.Now(),                                                        // Set the created_at timestamp
		UpdatedAt: time.Now(),                                                        // Set the updated_at timestamp
	}

	// Send session notification
	err = SendSessionNotification(notification) // Send the session notification
	if err != nil {                             // Check if there is an error sending the notification
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"}) // Return an error response
		return                                                                                // Return from the function
	} // Send session notification

	// Delete the session from the database
	err = DeleteSession(objectSessionID) // Delete the session
	if err != nil {                      // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session canceled and notification sent successfully"}) // Return a success response
}

// ArchiveSession: Archives a session and sends a notification to the creator
func ArchiveSession(c *gin.Context) { // Archive a session
	sessionID := c.Param("sessionId") // Get the session ID from the URL

	// Convert sessionID to ObjectID
	objectSessionID, err := primitive.ObjectIDFromHex(sessionID) // Convert the session ID to an ObjectID
	if err != nil {                                              // Check if there is an error converting the ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"}) // Return a bad request response
		return                                                              // Return from the function
	}

	// Find the session to get details
	var session models.Session
	var user models.User                                                                             // Define a session variable
	err = sessionCollection.FindOne(context.TODO(), bson.M{"_id": objectSessionID}).Decode(&session) // Find the session by ID
	if err != nil {                                                                                  // Check if there is an error
		if err == mongo.ErrNoDocuments { // Check if the session was not found
			c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"}) // Return a not found response
		} else { // If there's another error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		}
		return // Return from the function
	}

	// Prepare notification data
	notification := models.Notification{ // Define a notification variable with the required fields
		ID:        primitive.NewObjectID(),                                           // Generate a new ObjectID for the notification
		UserID:    user.ID,                                                           // Assuming session.CreatorID holds the ID of the user to notify
		Type:      "Session Archived",                                                // Set the notification type
		Message:   fmt.Sprintf("The session '%s' has been archived.", session.Title), // Message for the notification
		CreatedAt: time.Now(),                                                        // Set the created_at timestamp
		UpdatedAt: time.Now(),                                                        // Set the updated_at timestamp
	}

	// Send session notification
	err = SendSessionNotification(notification) // Send the session notification
	if err != nil {                             // Check if there is an error sending the notification
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"}) // Return an error response
		return                                                                                // Return from the function
	}

	// Update the session status to "archived"
	update := bson.M{"$set": bson.M{"status": "archived"}}                                       // Set the status to "archived"
	_, err = sessionCollection.UpdateOne(context.TODO(), bson.M{"_id": objectSessionID}, update) // Update the session
	if err != nil {                                                                              // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session archived and notification sent successfully"}) // Return a success response
}
