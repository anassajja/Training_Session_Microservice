package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the structure of a user document in MongoDB.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`            // Unique identifier for the user
	Name      string             `json:"name" bson:"name,omitempty"`         // Name of the user
	Email     string             `json:"email" bson:"email,omitempty"`       // Email address of the user
	Role      string             `json:"role" bson:"role,omitempty"`         // Role of the user (e.g., "admin", "user")
	Cin       string             `json:"cin" bson:"cin,omitempty"`           // National ID or CIN of the user
	Password  string             `json:"password" bson:"password,omitempty"` // Encrypted password of the user
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`        // Timestamp when the user was created
	UpdatedAt time.Time          `bson:"updatedAt" json:"updated_at"`        // Timestamp when the user was last updated
}
