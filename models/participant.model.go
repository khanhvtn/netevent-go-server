package models

import (
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Participant struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt            time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt            time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsValid              bool               `bson:"isValid" json:"isValid"`
	IsAttended           bool               `bson:"isAttended" json:"isAttended"`
	Event                Event              `bson:"event" json:"event"`
	Email                string             `bson:"email" json:"email"`
	Name                 string             `bson:"name" json:"name"`
	Academic             string             `bson:"academic" json:"academic"`
	School               string             `bson:"school" json:"school"`
	Major                string             `bson:"major" json:"major"`
	Phone                string             `bson:"phone" json:"phone"`
	DOB                  time.Time          `bson:"dob" json:"dob"`
	ExpectedGraduateDate time.Time          `bson:"expectedGraduateDate" json:"expectedGraduateDate"`
}

/* Input Type */
type NewParticipant struct {
	Event                string    `bson:"event" json:"event"`
	Email                string    `bson:"email" json:"email"`
	Name                 string    `bson:"name" json:"name"`
	Academic             string    `bson:"academic" json:"academic"`
	School               string    `bson:"school" json:"school"`
	Major                string    `bson:"major" json:"major"`
	Phone                string    `bson:"phone" json:"phone"`
	DOB                  time.Time `bson:"dob" json:"dob"`
	ExpectedGraduateDate time.Time `bson:"expectedGraduateDate" json:"expectedGraduateDate"`
}

/* Model Function */
/* GetAll: get all data based on condition*/
func (u *Participant) GetAll() ([]*Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("participants")
	defer cancel()

	//create an empty array to store all fields from collection
	var participants []*Participant = make([]*Participant, 0)

	//get all record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var participant Participant
		cur.Decode(&participant)
		participants = append(participants, &participant)
	}
	//response data to client
	if participants == nil {
		return make([]*Participant, 0), nil
	}
	return participants, nil
}

/*GetOne: get one record from a collection  */
func (u *Participant) GetOne(filter bson.M) (*Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("participant")
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

	participant := Participant{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&participant); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &participant, nil
}

/*Create: create a new record to a collection*/
func (u *Participant) Create(newParticipant NewParticipant) (*Participant, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("participants")
	defer cancel()

	//get event
	var event *Event
	if targetEvent, err := EventInstance.GetOne(bson.M{"_id": newParticipant.Event}); err != nil {
		return nil, err
	} else if targetEvent == nil {
		return nil, errors.New("invalid event id")
	} else {
		event = targetEvent
	}

	//convert to bson.M
	currentTime := time.Now()
	participant := Participant{
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		IsValid:              false,
		IsAttended:           false,
		Event:                *event,
		Email:                newParticipant.Email,
		Name:                 newParticipant.Name,
		Academic:             newParticipant.Academic,
		School:               newParticipant.School,
		Major:                newParticipant.Major,
		Phone:                newParticipant.Phone,
		DOB:                  newParticipant.DOB,
		ExpectedGraduateDate: newParticipant.ExpectedGraduateDate,
	}
	newData, err := utilities.InterfaceToBsonM(participant)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &Participant{
		ID:                   insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		IsValid:              false,
		IsAttended:           false,
		Event:                *event,
		Email:                newParticipant.Email,
		Name:                 newParticipant.Name,
		Academic:             newParticipant.Academic,
		School:               newParticipant.School,
		Major:                newParticipant.Major,
		Phone:                newParticipant.Phone,
		DOB:                  newParticipant.DOB,
		ExpectedGraduateDate: newParticipant.ExpectedGraduateDate,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u Participant) UpdateOne(filter bson.M, update bson.M) (*Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("participants")
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
	participant, errQuery := u.GetOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return participant, nil
}

//DeleteOne func is to update one record from a collection
func (u Participant) DeleteOne(filter bson.M) (*Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("participants")
	defer cancel()

	participant, errGet := u.GetOne(filter)
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

	return participant, nil
}

//export an instance Participant model
var ParticipantInstance = Participant{}
