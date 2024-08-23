package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pitch represents the structure of a pitch document in MongoDB.
type Pitch struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`        // Unique identifier for the pitch
	SessionID   primitive.ObjectID `bson:"session_id" json:"session_id"`   // ID of the associated training session
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`         // ID of the user who created the pitch
	Title       string             `bson:"title" json:"title"`             // Title of the pitch
	Description string             `bson:"description" json:"description"` // Description of the pitch
	CreatedAt   time.Time          `bson:"createdAt" json:"created_at"`    // Timestamp when the pitch was created
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updated_at"`    // Timestamp when the pitch was last updated
}
