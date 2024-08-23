package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session represents the structure of a session document in MongoDB.
type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`           // Unique identifier for the session
	Title        string             `bson:"title" json:"title"`                // Title of the session
	Description  string             `bson:"description" json:"description"`    // Description of the session
	StartTime    time.Time          `bson:"startTime" json:"start_time"`       // Start time of the session
	EndTime      time.Time          `bson:"endTime" json:"end_time"`           // End time of the session
	Location     string             `bson:"location" json:"location"`          // Location of the session
	TrainingType string             `bson:"trainingType" json:"training_type"` // Type of training
	Duration     int                `bson:"duration" json:"duration"`          // Duration of the session in minutes
	Recurrence   string             `bson:"recurrence" json:"recurrence"`      // Recurrence pattern of the session
	Coach        string             `bson:"coach" json:"coach"`                // Coach for the session
	CoachAssists []string           `bson:"coachAssists" json:"coach_assists"` // List of assistants for the coach
	Participants []string           `bson:"participants" json:"participants"`  // List of participants
	Status       string             `bson:"status" json:"status"`              // Status of the session (e.g., "active", "cancelled")
	QRCode       string             `bson:"qrCode" json:"qr_code"`             // QR code associated with the session
	CreatedAt    time.Time          `bson:"createdAt" json:"created_at"`      // Timestamp when the session was created
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updated_at"`       // Timestamp when the session was last updated
}
