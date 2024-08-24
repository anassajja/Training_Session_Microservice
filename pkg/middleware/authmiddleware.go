package middleware

import (
	"log"
	"net/http"
	"strings"

	"training_session/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var cfg *config.Config // Config variable to store the configuration

func init() { // Initialize the configuration
	cfg = config.LoadConfig() // Load the configuration
}

func AuthMiddleware() gin.HandlerFunc { // AuthMiddleware function to authenticate requests
	return func(c *gin.Context) { // Return a Gin handler function
		tokenString := c.GetHeader("Authorization") // Get the Authorization header from the request
		if tokenString == "" { // Check if the token is missing
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"}) // Return an error if the token is missing
			c.Abort() // Abort the request
			return // Return from the function
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ") // Remove the "Bearer " prefix from the token
		log.Printf("Token after prefix removal: %s", tokenString) // Log the token after removing the prefix

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // Parse the JWT token
			return []byte(cfg.JwtSecretKey), nil // Return the JWT secret key
		})

		if err != nil || !token.Valid { // Check if there is an error or the token is invalid
			log.Printf("Error parsing token: %v", err) // Log the error message
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"}) // Return an error response
			c.Abort() // Abort the request
			return   // Return from the function
		}

		c.Next() // Call the next handler
	}
}
