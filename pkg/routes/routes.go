package routes

import ( // Import the required packages
	"training_session/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) { // Set up the router

	// Add routes for users

	r.GET("/users", controllers.GetUsers)          // Add a route to get all users
	r.POST("/users", controllers.RegisterUser)     // Add a route to create a new user
	r.PUT("/users/:id", controllers.UpdateUser)    // Add a route to update a user
	r.DELETE("/users/:id", controllers.DeleteUser) // Add a route to delete a user

	// Refund Route
	// r.POST("/sessions/:userId/:sessionId/refund", controllers.ProcessRefund)

	// Add routes for sessions

	r.POST("/sessions", controllers.CreateSession)           // Add a route to create a new session
	r.PUT("/sessions/:id", controllers.UpdateSession)        // Add a route to update a session
	r.GET("/sessions", controllers.GetSessions)              // Add a route to get all sessions
	r.GET("/sessions/active", controllers.GetActiveSessions) // Add a route to get all active sessions
	r.GET("/sessions/:id", controllers.GetSessionByID)       // Add a route to get a session by ID
	r.GET("/sessions/user/:userId", controllers.GetSessionsByUserID)
	r.POST("/sessions/:userId/:sessionId/enroll", controllers.EnrollInSession)
	r.POST("/sessions/:sessionId/enrollments/:userId/cancel", controllers.CancelEnrollment)
	r.POST("/sessions/:id/cancel", controllers.CancelSession)
	r.POST("/sessions/:id/archive", controllers.ArchiveSession) // Route for archiving a session

	// Add routes for invitations
	r.POST("/invitations", controllers.SendInvitation)                // Add a route to send an invitation
	r.POST("/invitations/:id/accept", controllers.AcceptInvitation)   // Add a route to accept an invitation
	r.POST("/invitations/:id/decline", controllers.DeclineInvitation) // Add a route to decline an invitation
	r.GET("/invitations", controllers.GetInvitations)                 // Add a route to get all invitations
	r.GET("/invitations/:id", controllers.GetInvitationByID)          // Add a route to get an invitation by ID
	r.DELETE("/invitations/:id", controllers.DeleteInvitation)        // Add a route to delete an invitation

	// Add routes for QR codes
	r.GET("/sessions/:sessionId/qrcode", controllers.GenerateQRCode)
	r.GET("/sessions/:sessionId/validate", controllers.ValidateQRCode)

	// Add routes for Feedback routes
	r.POST("/feedback", controllers.SubmitFeedback)           // Route to submit feedback
	r.GET("/feedback/user/:userId", controllers.ViewFeedback) // Route to view feedback submitted by the user
	r.PUT("/feedback/:id", controllers.EditFeedback)          // Route to edit feedback by ID
	r.DELETE("/feedback/:id", controllers.DeleteFeedback)     // Route to delete feedback by ID

	// Add routes for Coaches
	/* 	r.GET("/coaches", controllers.GetCoaches)         // Add a route to get all coaches
	   	r.POST("/coaches", controllers.RegisterCoach)     // Add a route to create a new coach
	   	r.PUT("/coaches/:id", controllers.UpdateCoach)    // Add a route to update a coach
	   	r.DELETE("/coaches/:id", controllers.DeleteCoach) // Add a route to delete a coach */

	// Add routes for Notifications
	r.POST("/notifications/user", controllers.SendUserNotification) // Send user notifications
	r.GET("/notifications/:userId", controllers.GetNotifications)   // Get all notifications for a user
	r.DELETE("/notifications/:id", controllers.DeleteNotification)  // Delete a notification by ID

	// Add routes for Pitch bookings
	r.POST("/pitches", controllers.BookPitch)                    // Add a route to book a pitch
	r.GET("/pitches", controllers.GetPitchBookings)              // Add a route to get all pitches
	r.GET("/pitches/:id", controllers.GetPitchBookingsByPitchID) // Add a route to get a pitch by pitch id
	r.GET("/pitches/:id", controllers.GetPitchBookingsByUserID)  // Add a route to get a pitch by user id
	r.PUT("/pitches/:id", controllers.UpdatePitchBooking)        // Add a route to update a pitch booking
	r.DELETE("/pitches/:id", controllers.DeletePitchBooking)     // Add a route to delete a pitch booking

} // End of the SetupRoutes function
