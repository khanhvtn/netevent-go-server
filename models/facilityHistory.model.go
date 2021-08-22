package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CollectionFacilityHistoryName = "facilityHistories"

type FacilityHistory struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
	Facility   primitive.ObjectID `bson:"facility" json:"facility"`
	BorrowDate time.Time          `bson:"borrowDate" json:"borrowDate"`
	ReturnDate time.Time          `bson:"returnDate" json:"returnDate"`
	Event      primitive.ObjectID `bson:"event,omitempty" json:"event"`
}
