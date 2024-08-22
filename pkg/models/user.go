package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct { // User model
	ID        primitive.ObjectID `bson:"_id,omitempty"`                      // Define an ID field
	Name      string             `json:"name" bson:"name,omitempty"`         // Name field
	Email     string             `json:"email" bson:"email,omitempty"`       // Email field
	Role      string             `json:"role" bson:"role,omitempty"`         // Role field
	Cin       string             `json:"cin" bson:"cin,omitempty"`           // Cin field
	Password  string             `json:"password" bson:"password,omitempty"` // Password field
	CreatedAt time.Time          `bson:"createdAt"`                          // Define a CreatedAt field
	UpdatedAt time.Time          `bson:"updatedAt"`                          // Define an UpdatedAt field
}
