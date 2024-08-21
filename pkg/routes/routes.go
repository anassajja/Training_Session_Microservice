package routes

import (
	"training_session/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    controllers.SetupUserCollection() // Initialize user collection here

    r.GET("/users", controllers.GetUsers)
    r.POST("/users", controllers.CreateUser)
}
