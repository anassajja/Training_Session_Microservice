package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct { // Define a Session struct
	ID           primitive.ObjectID `bson:"_id,omitempty"` // Define an ID field
	Title        string             `bson:"title"` 	   // Define a Title field
	Description  string             `bson:"description"`  // Define a Description field
	StartTime    time.Time          `bson:"startTime"`   // Define a StartTime field
	EndTime      time.Time          `bson:"endTime"`    // Define an EndTime field
	Location     string             `bson:"location"`  // Define a Location field
	TrainingType string             `bson:"trainingType"` // Define a TrainingType field
	Duration     int                `bson:"duration"` // Define a Duration field
	Recurrence   string             `bson:"recurrence"` // Define a Recurrence field
	Coach        string             `bson:"coach"` // Define a Coach field
	CoachAssists []string           `bson:"coachAssists"` // Define a CoachAssists field
	Participants []string           `bson:"participants"` // Define a Participants field
	Status       string             `bson:"status"` // Add Status field to represent the session status (e.g., "active", "cancelled")
	CreatedAt    time.Time          `bson:"createdAt"` // Define a CreatedAt field
	UpdatedAt    time.Time          `bson:"updatedAt"` // Define an UpdatedAt field
}
