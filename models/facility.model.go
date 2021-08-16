package models

import (
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

/* Input Type */
type NewFacility struct {
	Name string `bson:"name" json:"name"`
	Code string `bson:"code" json:"code"`
	Type string `bson:"type" json:"type"`
}

/* Model Function */
/* GetAll: get all data based on condition*/
func (u *Facility) GetAll() ([]*Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilities")
	defer cancel()

	//create an empty array to store all fields from collection
	var facilities []*Facility = make([]*Facility, 0)

	//get all record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var facility Facility
		cur.Decode(&facility)
		facilities = append(facilities, &facility)
	}
	//response data to client
	if facilities == nil {
		return make([]*Facility, 0), nil
	}
	return facilities, nil
}

/*GetOne: get one record from a collection  */
func (u *Facility) GetOne(filter bson.M) (*Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilities")
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

	facility := Facility{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&facility); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &facility, nil
}

/*Create: create a new record to a collection*/
func (u *Facility) Create(newFacility NewFacility) (*Facility, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilities")
	defer cancel()

	//convert to bson.M
	currentTime := time.Now()
	facility := Facility{
		Name:      newFacility.Name,
		Code:      newFacility.Code,
		Type:      newFacility.Type,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		IsDeleted: false,
		Status:    false,
	}
	newData, err := utilities.InterfaceToBsonM(facility)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &Facility{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      newFacility.Name,
		Code:      newFacility.Code,
		Type:      newFacility.Type,
		IsDeleted: false,
		Status:    false,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u Facility) UpdateOne(filter bson.M, update bson.M) (*Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilities")
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
	facility, errQuery := u.GetOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return facility, nil
}

//DeleteOne func is to update one record from a collection
func (u Facility) DeleteOne(filter bson.M) (*Facility, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("facilities")
	defer cancel()

	facility, errGet := u.GetOne(filter)
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

	return facility, nil
}

//export an instance Facility model
var FacilityInstance = Facility{}
