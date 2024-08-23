package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification represents the structure of a notification document in MongoDB.
type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`     // Unique identifier for the notification
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`      // ID of the user associated with the notification
	Type      string             `bson:"type" json:"type"`            // Type of the notification (e.g., "User", "Session")
	Message   string             `bson:"message" json:"message"`      // Content of the notification
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"` // Timestamp when the notification was created
	UpdatedAt time.Time          `bson:"updatedAt" json:"updated_at"` // Timestamp when the notification was last updated
}
