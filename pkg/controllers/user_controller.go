package controllers

import (
	"context"
	"net/http"
	"training_session/db"
	"training_session/pkg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
    var users []models.User // Define a users variable
    cursor, err := userCollection.Find(context.TODO(), bson.D{{}}) // Find all users
    if err != nil { // Check if there is an error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
        return // Return from the function
    }
    defer cursor.Close(context.TODO()) // Close the cursor
    for cursor.Next(context.TODO()) { // Iterate over the cursor
        var user models.User // Define a user variable
        if err := cursor.Decode(&user); err != nil { // Decode the user
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // Return an error response
            return // Return from the function
        }
        users = append(users, user) // Append the user to the users slice
    }
    c.JSON(http.StatusOK, users) // Return a success response
}

func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    _, err := userCollection.InsertOne(context.TODO(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}
