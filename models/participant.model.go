package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CollectionParticipantName = "participants"

type Participant struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt            time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsValid              bool               `bson:"isValid" json:"isValid"`
	IsAttended           bool               `bson:"isAttended" json:"isAttended"`
	Event                primitive.ObjectID `bson:"event" json:"event"`
	Email                string             `bson:"email" json:"email"`
	Name                 string             `bson:"name" json:"name"`
	Academic             string             `bson:"academic" json:"academic"`
	School               string             `bson:"school" json:"school"`
	Major                string             `bson:"major" json:"major"`
	Phone                string             `bson:"phone" json:"phone"`
	DOB                  time.Time          `bson:"dob" json:"dob"`
	ExpectedGraduateDate time.Time          `bson:"expectedGraduateDate" json:"expectedGraduateDate"`
}
