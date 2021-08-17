package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Facility struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Status    bool               `bson:"status" json:"status"`
	Name      string             `bson:"name" json:"name"`
	Code      string             `bson:"code" json:"code"`
	Type      string             `bson:"type" json:"type"`
	IsDeleted bool               `bson:"isDeleted" json:"isDeleted"`
}
