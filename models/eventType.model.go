package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Name      string             `bson:"name" json:"name"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted"`
}
