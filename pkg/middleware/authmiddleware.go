package middleware

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var jwtSecret []byte // Define a global variable to store the JWT secret key

func init() { // Initialize the environment variables and JWT secret key when the package is imported
	// Load environment variables from .env file
	err := godotenv.Load() // Load the .env file
	if err != nil {        // Check if there is an error
		panic("Error loading .env file") // Panic if there is an error loading the .env file to stop the application
	}

	// Get the secret key from the environment variables
	jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func AuthMiddleware() gin.HandlerFunc { // AuthMiddleware function to authenticate requests using JWT tokens as middleware for protected routes (e.g., sessions, feedback, notifications)
	return func(c *gin.Context) { // Return a function that handles the request
		tokenString := c.GetHeader("Authorization") // Get the token from the Authorization header
		if tokenString == "" {                      // Check if the token is missing
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"}) // Return an error response
			c.Abort()                                                        // Abort the request
			return                                                           // Abort the request
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // Parse the token
			return jwtSecret, nil // Return the secret key
		})

		if err != nil || !token.Valid { // Check if there is an error or the token is invalid
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"}) // Return an error response
			c.Abort()                                                        // Abort the request
			return                                                           // Abort the request
		}

		c.Next() // Continue to the next middleware
	}
}
