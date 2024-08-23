package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feedback model
type Feedback struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // Unique identifier for the feedback
	SessionID primitive.ObjectID `bson:"session_id" json:"session_id"` // Session for which the feedback was provided
	CoachID   primitive.ObjectID `bson:"coach_id" json:"coach_id"` 		// Coach who received the feedback
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"` // User who provided the feedback
	Content   string             `bson:"content" json:"content"` // Feedback content
	Rating    int                `bson:"rating" json:"rating"` // Example: Rating between 1 to 5
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // Timestamp when the feedback was created
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"` // Timestamp when the feedback was last updated
}
