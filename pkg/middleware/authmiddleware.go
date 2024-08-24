package middleware

import (
	"log"
	"net/http"
	"strings"

	"training_session/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		log.Printf("Token after prefix removal: %s", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Error parsing token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
