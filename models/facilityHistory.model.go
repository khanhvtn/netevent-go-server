package models

import (
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FacilityHistory struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
	Facility   Facility           `bson:"facility" json:"facility"`
	BorrowDate time.Time          `bson:"borrowDate" json:"borrowDate"`
	ReturnDate time.Time          `bson:"returnDate" json:"returnDate"`
	Event      Event              `bson:"event" json:"event"`
}

/* Input Type */
type NewFacilityHistory struct {
	Facility   string    `bson:"facility" json:"facility"`
	BorrowDate time.Time `bson:"borrowDate" json:"borrowDate"`
	ReturnDate time.Time `bson:"returnDate" json:"returnDate"`
	Event      string    `bson:"event" json:"event"`
}

/* Model Function */
/* GetAll: get all data based on condition*/
func (u *FacilityHistory) GetAll() ([]*FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilityHistories")
	defer cancel()

	//create an empty array to store all fields from collection
	var facilityHistories []*FacilityHistory = make([]*FacilityHistory, 0)

	//get all record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var facilityHistory FacilityHistory
		cur.Decode(&facilityHistory)
		facilityHistories = append(facilityHistories, &facilityHistory)
	}
	//response data to client
	if facilityHistories == nil {
		return make([]*FacilityHistory, 0), nil
	}
	return facilityHistories, nil
}

/*GetOne: get one record from a collection  */
func (u *FacilityHistory) GetOne(filter bson.M) (*FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilityHistories")
	defer cancel()

	//convert id to object id when filter contain _id
	if checkID := filter["_id"]; checkID != nil {
		if _, ok := checkID.(primitive.ObjectID); !ok {
			id, err := primitive.ObjectIDFromHex(checkID.(string))
			if err != nil {
				return nil, err
			}
			filter["_id"] = id
		}
	}

	facilityHistory := FacilityHistory{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&facilityHistory); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &facilityHistory, nil
}

/*Create: create a new record to a collection*/
func (u *FacilityHistory) Create(newFacilityHistory NewFacilityHistory) (*FacilityHistory, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilityHistories")
	defer cancel()

	//get facility, event
	var facility *Facility
	var event *Event
	if targetFacility, err := FacilityInstance.GetOne(bson.M{"_id": newFacilityHistory.Facility}); err != nil {
		return nil, err
	} else if targetFacility == nil {
		return nil, errors.New("invalid facility id")
	} else {
		facility = targetFacility
	}

	if targetEvent, err := EventInstance.GetOne(bson.M{"_id": newFacilityHistory.Event}); err != nil {
		return nil, err
	} else if targetEvent == nil {
		return nil, errors.New("invalid event id")
	} else {
		event = targetEvent
	}

	//convert to bson.M
	currentTime := time.Now()
	facilityHistory := FacilityHistory{
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
		Facility:   *facility,
		BorrowDate: newFacilityHistory.BorrowDate,
		ReturnDate: newFacilityHistory.ReturnDate,
		Event:      *event,
	}
	newData, err := utilities.InterfaceToBsonM(facilityHistory)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &FacilityHistory{
		ID:         insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
		Facility:   *facility,
		BorrowDate: newFacilityHistory.BorrowDate,
		ReturnDate: newFacilityHistory.ReturnDate,
		Event:      *event,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u FacilityHistory) UpdateOne(filter bson.M, update bson.M) (*FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilityHistories")
	defer cancel()

	//convert id to object id when filter contain _id
	if checkID := filter["_id"]; checkID != nil {
		if _, ok := checkID.(primitive.ObjectID); !ok {
			id, err := primitive.ObjectIDFromHex(checkID.(string))
			if err != nil {
				return nil, err
			}
			filter["_id"] = id
		}
	}

	//update user information
	newUpdate := bson.M{"$set": update}
	updateResult, err := collection.UpdateOne(ctx, filter, newUpdate)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, errors.New("id not found")
	}

	//query the new update
	facilityHistory, errQuery := u.GetOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return facilityHistory, nil
}

//DeleteOne func is to update one record from a collection
func (u FacilityHistory) DeleteOne(filter bson.M) (*FacilityHistory, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilityHistories")
	defer cancel()

	facilityHistory, errGet := u.GetOne(filter)
	if errGet != nil {
		return nil, errGet
	}

	//delete user from database
	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		//response to client if there is an error.
		return nil, err
	}

	if deleteResult.DeletedCount == 0 {
		return nil, errors.New("id not found")
	}

	return facilityHistory, nil
}

//export an instance FacilityHistory model
var FacilityHistoryInstance = FacilityHistory{}
