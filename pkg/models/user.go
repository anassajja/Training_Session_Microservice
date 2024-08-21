package models

type User struct { // User model
	ID    string `json:"id" bson:"_id,omitempty"`      // ID field
	Name  string `json:"name" bson:"name,omitempty"`   // Name field
	Email string `json:"email" bson:"email,omitempty"` // Email field
}
