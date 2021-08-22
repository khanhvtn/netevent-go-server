package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CollectionUserName = "users"

/* Model Type */
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password,omitempty"`
	Roles     []string           `bson:"roles" json:"roles"`
}
