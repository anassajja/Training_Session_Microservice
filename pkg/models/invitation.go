package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Invitation represents the structure of an invitation document in MongoDB.
type Invitation struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`                // Unique identifier for the invitation
	SessionID      string             `bson:"session_id" json:"session_id"`           // ID of the associated training session
	UserID         string             `bson:"user_id" json:"user_id"`                 // ID of the user who receives the invitation
	InvitationDate string             `bson:"invitation_date" json:"invitation_date"` // Date when the invitation was sent
	Status         string             `bson:"status" json:"status"`                   // Status of the invitation (e.g., "pending", "accepted", "declined")
	CreatedAt      time.Time          `bson:"createdAt" json:"created_at"`            // Timestamp when the session was created
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updated_at"`            // Timestamp when the session was last updated
}
