package routes

import ( // Import the required packages
	"training_session/pkg/controllers"
	"training_session/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) { // SetupRoutes function to define the routes
	// Add routes for users
	r.GET("/users", controllers.GetUsers)                     // Define a route to get all users
	r.GET("/users/:userId", controllers.GetUserByID)          // Define a route to get a user by ID
	r.POST("/users/register", controllers.RegisterUser)       // Define a route to register a new user
	r.POST("/users/login", controllers.LoginUser)             // Define a route to login a user
	r.POST("/users/logout", controllers.LogoutUser)           // Define a route to logout a user
	r.PUT("/users/update/:userId", controllers.UpdateUser)    // Define a route to update a user
	r.DELETE("/users/delete/:userId", controllers.DeleteUser) // Define a route to delete a user

	// Protected routes with authentication middleware
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Use the AuthMiddleware to authenticate requests

	// Add routes for sessions
	protected.POST("/sessions", controllers.CreateSession)                                              // Define a route to create a new session
	protected.PUT("/sessions/:sessionId", controllers.UpdateSession)                                    // Define a route to update a session
	protected.GET("/sessions", controllers.GetSessions)                                                 // Define a route to get all sessions
	protected.GET("/sessions/active", controllers.GetActiveSessions)                                    // Define a route to get all active sessions
	protected.GET("/sessions/:sessionId", controllers.GetSessionByID)                                   // Define a route to get a session by ID
	protected.GET("/sessions/user/:userId", controllers.GetSessionsByUserID)                            // Define a route to get all sessions by a user ID
	protected.POST("/sessions/:sessionId/user/:userId/enroll", controllers.EnrollInSession)             // Define a route to enroll in a session
	protected.POST("/sessions/:sessionId/user/:userId/cancel-enrollment", controllers.CancelEnrollment) // Define a route to cancel enrollment in a session
	protected.POST("/sessions/:sessionId/cancel", controllers.CancelSession)                            // Define a route to cancel a session
	protected.POST("/sessions/:sessionId/archive", controllers.ArchiveSession)                          // Define a route to archive a session

	// Add routes for invitations
	r.POST("/invitations", controllers.SendInvitation)                          // Define a route to send an invitation
	r.POST("/invitations/:invitationId/accept", controllers.AcceptInvitation)   // Define a route to accept an invitation
	r.POST("/invitations/:invitationId/decline", controllers.DeclineInvitation) // Define a route to decline an invitation
	r.GET("/invitations", controllers.GetInvitations)                           // Define a route to get all invitations
	r.GET("/invitations/:invitationId", controllers.GetInvitationByID)          // Define a route to get an invitation by ID
	r.DELETE("/invitations/:invitationId", controllers.DeleteInvitation)        // Define a route to delete an invitation

	// Add routes for QR codes
	r.GET("/sessions/:sessionId/qrcode", controllers.GenerateQRCode)   // Define a route to generate a QR code
	r.GET("/sessions/:sessionId/validate", controllers.ValidateQRCode) // Define a route to validate a QR code

	// Add routes for Feedback
	r.POST("/feedback", controllers.SubmitFeedback)               // Define a route to submit feedback
	r.GET("/feedback/user/:userId", controllers.ViewFeedback)     // Define a route to view feedback
	r.PUT("/feedback/:feedbackId", controllers.EditFeedback)      // Define a route to edit feedback
	r.DELETE("/feedback/:feedbackId", controllers.DeleteFeedback) // Define a route to delete feedback

	// Add routes for Notifications
	r.POST("/notifications/user", controllers.SendUserNotification)            // Define a route to send a user notification
	r.GET("/notifications/:userId", controllers.GetNotifications)              // Define a route to get notifications for a user
	r.DELETE("/notifications/:notificationId", controllers.DeleteNotification) // Define a route to delete a notification

	// Add routes for Pitch bookings
	r.POST("/pitches", controllers.BookPitch)                            // Define a route to book a pitch
	r.GET("/pitches", controllers.GetPitchBookings)                      // Define a route to get all pitch bookings
	r.GET("/pitches/:pitchId", controllers.GetPitchBookingsByPitchID)    // Define a route to get pitch bookings by pitch ID
	r.GET("/pitches/user/:userId", controllers.GetPitchBookingsByUserID) // Define a route to get pitch bookings by user ID
	r.PUT("/pitches/:pitchId", controllers.UpdatePitchBooking)           // Define a route to update a pitch booking
	r.DELETE("/pitches/:pitchId", controllers.DeletePitchBooking)        // Define a route to delete a pitch booking
} // End of SetupRoutes function
