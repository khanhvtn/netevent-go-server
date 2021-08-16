package models

import (
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/* Model Type */
type Event struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt             time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt             time.Time          `bson:"updatedAt" json:"updatedAt"`
	Tags                  []string           `bson:"tags" json:"tags"`
	IsApproved            bool               `bson:"isApproved" json:"isApproved"`
	Reviewer              User               `bson:"reviewer;omitempty" json:"reviewer"`
	IsFinished            bool               `bson:"isFinished" json:"isFinished"`
	Tasks                 []Task             `bson:"tasks" json:"tasks"`
	FacilityHistories     []FacilityHistory  `bson:"facilityHistories" json:"facilityHistories"`
	Name                  string             `bson:"name" json:"name"`
	Language              string             `bson:"language" json:"language"`
	EventType             EventType          `bson:"eventType" json:"eventType"`
	Mode                  string             `bson:"mode" json:"mode"`
	Location              string             `bson:"location" json:"location"`
	Accommodation         string             `bson:"accommodation" json:"accommodation"`
	RegistrationCloseDate time.Time          `bson:"registrationCloseDate" json:"registrationCloseDate"`
	StartDate             time.Time          `bson:"startDate" json:"startDate"`
	EndDate               time.Time          `bson:"endDate" json:"endDate"`
	MaxParticipants       int                `bson:"maxParticipants" json:"maxParticipants"`
	Description           string             `bson:"description" json:"description"`
	Owner                 User               `bson:"owner" json:"owner"`
	Budget                float64            `bson:"budget" json:"budget"`
	Image                 string             `bson:"image" json:"image"`
	IsDeleted             bool               `bson:"isDeleted" json:"isDeleted"`
}

/* Input Type */
type NewEvent struct {
	Tags                  []string          `bson:"tags" json:"tags"`
	Reviewer              string            `bson:"reviewer" json:"reviewer"`
	Tasks                 []Task            `bson:"tasks" json:"tasks"`
	FacilityHistories     []FacilityHistory `bson:"facilityHistories" json:"facilityHistories"`
	Name                  string            `bson:"name" json:"name"`
	Language              string            `bson:"language" json:"language"`
	EventType             string            `bson:"eventType" json:"eventType"`
	Mode                  string            `bson:"mode" json:"mode"`
	Location              string            `bson:"location" json:"location"`
	Accommodation         string            `bson:"accommodation" json:"accommodation"`
	RegistrationCloseDate time.Time         `bson:"registrationCloseDate" json:"registrationCloseDate"`
	StartDate             time.Time         `bson:"startDate" json:"startDate"`
	EndDate               time.Time         `bson:"endDate" json:"endDate"`
	MaxParticipants       int               `bson:"maxParticipants" json:"maxParticipants"`
	Description           string            `bson:"description" json:"description"`
	Owner                 string            `bson:"owner" json:"owner"`
	Budget                float64           `bson:"budget" json:"budget"`
	Image                 string            `bson:"image" json:"image"`
}

/* Model Function */
/* GetAll: get all data based on condition*/
func (u *Event) GetAll() ([]*Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("events")
	defer cancel()

	//create an empty array to store all fields from collection
	var events []*Event = make([]*Event, 0)

	//get all record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var event Event
		cur.Decode(&event)
		events = append(events, &event)
	}
	//response data to client
	if events == nil {
		return make([]*Event, 0), nil
	}
	return events, nil
}

/*GetOne: get one record from a collection  */
func (u *Event) GetOne(filter bson.M) (*Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("events")
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

	event := Event{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&event); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &event, nil
}

/*Create: create a new record to a collection*/
func (u *Event) Create(newEvent NewEvent) (*Event, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("events")
	defer cancel()

	//get reviewer, eventType, owner
	var reviewer, owner *User
	var eventType *EventType
	if targetReviewer, err := UserInstance.GetOne(bson.M{"_id": newEvent.Reviewer}); err != nil {
		return nil, err
	} else if targetReviewer == nil {
		return nil, errors.New("invalid reviewer id")
	} else {
		reviewer = targetReviewer
	}

	if targetOwner, err := UserInstance.GetOne(bson.M{"_id": newEvent.Owner}); err != nil {
		return nil, err
	} else if targetOwner == nil {
		return nil, errors.New("invalid owner id")
	} else {
		owner = targetOwner
	}
	if targetEventType, err := EventTypeInstance.GetOne(bson.M{"_id": newEvent.EventType}); err != nil {
		return nil, err
	} else if targetEventType == nil {
		return nil, errors.New("invalid eventType id")
	} else {
		eventType = targetEventType
	}
	//convert o bson.M
	currentTime := time.Now()
	event := Event{
		Tags:                  newEvent.Tags,
		IsApproved:            false,
		Reviewer:              *reviewer,
		IsFinished:            false,
		Tasks:                 newEvent.Tasks,
		FacilityHistories:     newEvent.FacilityHistories,
		Name:                  newEvent.Name,
		Language:              newEvent.Language,
		EventType:             *eventType,
		Mode:                  newEvent.Mode,
		Location:              newEvent.Location,
		Accommodation:         newEvent.Accommodation,
		RegistrationCloseDate: newEvent.RegistrationCloseDate,
		StartDate:             newEvent.StartDate,
		EndDate:               newEvent.EndDate,
		MaxParticipants:       newEvent.MaxParticipants,
		Description:           newEvent.Description,
		Owner:                 *owner,
		Budget:                newEvent.Budget,
		Image:                 newEvent.Image,
		IsDeleted:             false,
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
	}
	newData, err := utilities.InterfaceToBsonM(event)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:                    insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
		Tags:                  newEvent.Tags,
		IsApproved:            false,
		Reviewer:              *reviewer,
		IsFinished:            false,
		Tasks:                 newEvent.Tasks,
		FacilityHistories:     newEvent.FacilityHistories,
		Name:                  newEvent.Name,
		Language:              newEvent.Language,
		EventType:             *eventType,
		Mode:                  newEvent.Mode,
		Location:              newEvent.Location,
		Accommodation:         newEvent.Accommodation,
		RegistrationCloseDate: newEvent.RegistrationCloseDate,
		StartDate:             newEvent.StartDate,
		EndDate:               newEvent.EndDate,
		MaxParticipants:       newEvent.MaxParticipants,
		Description:           newEvent.Description,
		Owner:                 *owner,
		Budget:                newEvent.Budget,
		Image:                 newEvent.Image,
		IsDeleted:             false}, nil
}

/*UpdateOne: update one record from a collection*/
func (u Event) UpdateOne(filter bson.M, update bson.M) (*Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("events")
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
	event, errQuery := u.GetOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return event, nil
}

//DeleteOne func is to update one record from a collection
func (u Event) DeleteOne(filter bson.M) (*Event, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("events")
	defer cancel()

	event, errGet := u.GetOne(filter)
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

	return event, nil
}

/* Input Functions */

//export an instance User model
var EventInstance = Event{}
