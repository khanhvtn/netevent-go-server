package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CollectionEventName = "events"

/* Model Type */
type Event struct {
	ID                    primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	CreatedAt             time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt             time.Time            `bson:"updatedAt" json:"updatedAt"`
	Tags                  []string             `bson:"tags" json:"tags"`
	IsApproved            bool                 `bson:"isApproved" json:"isApproved"`
	Reviewer              *primitive.ObjectID  `bson:"reviewer" json:"reviewer"`
	IsFinished            bool                 `bson:"isFinished" json:"isFinished"`
	Tasks                 []primitive.ObjectID `bson:"tasks,omitempty" json:"tasks"`
	FacilityHistories     []primitive.ObjectID `bson:"facilityHistories,omitempty" json:"facilityHistories"`
	Name                  string               `bson:"name" json:"name"`
	Language              string               `bson:"language" json:"language"`
	EventType             primitive.ObjectID   `bson:"eventType" json:"eventType"`
	Mode                  string               `bson:"mode" json:"mode"`
	Location              string               `bson:"location" json:"location"`
	Accommodation         string               `bson:"accommodation" json:"accommodation"`
	RegistrationCloseDate time.Time            `bson:"registrationCloseDate" json:"registrationCloseDate"`
	StartDate             time.Time            `bson:"startDate" json:"startDate"`
	EndDate               time.Time            `bson:"endDate" json:"endDate"`
	MaxParticipants       int                  `bson:"maxParticipants" json:"maxParticipants"`
	Description           string               `bson:"description" json:"description"`
	Owner                 primitive.ObjectID   `bson:"owner" json:"owner"`
	Budget                float64              `bson:"budget" json:"budget"`
	Image                 string               `bson:"image" json:"image"`
	IsDeleted             bool                 `bson:"isDeleted" json:"isDeleted"`
}
