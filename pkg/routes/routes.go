package routes

import ( // Import the required packages
	"training_session/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) { // Set up the router

	// Add routes for users

	r.GET("/users", controllers.GetUsers)          // Add a route to get all users
	r.POST("/users", controllers.CreateUser)       // Add a route to create a new user
	r.PUT("/users/:id", controllers.UpdateUser)    // Add a route to update a user
	r.DELETE("/users/:id", controllers.DeleteUser) // Add a route to delete a user

	// Add routes for sessions

	r.POST("/sessions", controllers.CreateSession)    // Add a route to create a new session
	r.PUT("/sessions/:id", controllers.UpdateSession) // Add a route to update a session

	r.GET("/sessions", controllers.GetSessions)        // Add a route to get all sessions
	r.GET("/sessions/:id", controllers.GetSessionByID) // Add a route to get a session by ID

	r.DELETE("/sessions/:id", controllers.DeleteSession) // Add a route to delete a session

} // End of the SetupRoutes function
