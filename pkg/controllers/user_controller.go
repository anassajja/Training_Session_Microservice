package controllers

import (
	"context"
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
	"golang.org/x/crypto/bcrypt"
)

var (
	userCollection *mongo.Collection // Define a userCollection variable
	cfg            *config.Config    // Define a config variable
)

func InitializeUser(database *mongo.Database) { // Initialize the controllers
	userCollection = database.Collection("users")       // Set the user collection
	sessionCollection = database.Collection("sessions") // Set the session collection
}

func init() { // Initialize the configuration
	cfg = config.LoadConfig() // Load the configuration
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // Hash the password
	if err != nil {                                                                               // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"}) // Return an error response
		return                                                                            // Return from the function
	}
	user.Password = string(hashedPassword) // Store the hashed password

	user.ID = primitive.NewObjectID() // Generate a new ObjectID for the user
	user.CreatedAt = time.Now()       // Set the created_at timestamp
	user.UpdatedAt = time.Now()       // Set the updated_at timestamp

	_, err = userCollection.InsertOne(context.TODO(), user) // Insert the user
	if err != nil {                                         // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
		return                                                              // Return from the function
	}
	c.JSON(http.StatusCreated, user) // Return the created user
}

func LoginUser(c *gin.Context) { // Login a user
	var user models.User // Define a user variable
	if err := c.ShouldBindJSON(&user); err != nil { // Bind the JSON to the user struct
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"}) // Return an error response
		return
	}

	// Retrieve the user from the database
	var foundUser models.User                                                                     // Define a user variable
	err := userCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&foundUser) // Find the user by email
	if err != nil {                                                                               // Check if there is an error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"}) // Return an error response
		return                                                            // Return from the function
	}

	// Check if the password is correct (use proper hashing in production)
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)) // Compare the hashed password
	if err != nil {                                                                        // Check if there is an error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // Return an error response
		return                                                                 // Return from the function
	}

	// Create JWT token with the correct user ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // Create a new JWT token
		"id":  foundUser.ID.Hex(),                    // Convert ObjectID to hex string
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Set the expiration time to 72 hours
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey)) // Sign the token with the secret key
	if err != nil {                                                  // Check if there is an error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"}) // Return an error response
		return
	}

	c.SetCookie("auth_token", tokenString, int(time.Hour*2), "/", "", false, true) // Set the auth token cookie
	c.JSON(http.StatusOK, gin.H{                                                   // Return the user and token
		"user":  foundUser,   // Return the user
		"token": tokenString, // Return the token
	})
}

func LogoutUser(c *gin.Context) { // Logout a user
	// Clear the auth token cookie by setting it to an empty value and setting the expiration to -1
	c.SetCookie("auth_token", "", -1, "/", "", false, true) // Clear the auth token cookie

	// Log the logout action for debugging and auditing purposes
	userID, _ := c.Get("userID") // Retrieve the user ID from context (assuming it's set by authentication middleware)
	if userID != nil {           // Check if the user ID is not nil
		log.Printf("User %s logged out successfully", userID) // Log the user ID
	} else {
		log.Println("Anonymous user logged out") // Log a message for anonymous users
	}

	// Return a success response with a message
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"}) // Return a success response
}

func GetUserByID(c *gin.Context) { // Get a user by ID
	userID := c.Param("userId") // Get the user ID from the URL

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
	userID := c.Param("userId") // Get the user ID from the URL
	var user models.User        // Define a user variable

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

	// Hash the password if it is being updated
	if user.Password != "" { // Check if the password is not empty
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // Hash the password
		if err != nil {                                                                               // Check if there is an error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"}) // Return an error response
			return
		}
		user.Password = string(hashedPassword) // Store the hashed password
	}

	// Update the user document
	user.ID = objectID          // Set the user ID
	user.UpdatedAt = time.Now() // Set the updated_at timestamp

	filter := bson.M{"_id": objectID}                                       // Define the filter to find the user by ID
	update := bson.M{"$set": user}                                          // Define the update operation with the new user data (set)
	result, err := userCollection.UpdateOne(context.TODO(), filter, update) // Update the user document with the new data
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
	userID := c.Param("userId") // Get the user ID from the URL

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

/*
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
*/
