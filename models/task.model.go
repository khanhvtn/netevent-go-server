package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Event     primitive.ObjectID `bson:"event" json:"event"`
	Name      string             `bson:"name" json:"name"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	Type      string             `bson:"type" json:"type"`
	StartDate time.Time          `bson:"startDate" json:"startDate"`
	EndDate   time.Time          `bson:"endDate" json:"endDate"`
}
