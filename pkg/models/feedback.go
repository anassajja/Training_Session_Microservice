package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Feedback model
type Feedback struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SessionID primitive.ObjectID `bson:"session_id" json:"session_id"`
	CoachID   primitive.ObjectID `bson:"coach_id" json:"coach_id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Content   string             `bson:"content" json:"content"`
	Rating    int                `bson:"rating" json:"rating"` // Example: Rating between 1 to 5
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
