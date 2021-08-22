package services

import (
	"context"
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/database"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ParticipantRepositoryName = "ParticipantRepositoryName"

type ParticipantRepository struct {
	MongoCN *database.MongoInstance
}

/* createContextAndTargetCol: create and return targeted collection based on collection name */
func (u *ParticipantRepository) createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = u.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* FindAll: get all data based on condition*/
func (u *ParticipantRepository) FindAll(condition *bson.M) ([]*models.Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionParticipantName)
	defer cancel()

	//create an empty array to store all fields from collection
	var participants []*models.Participant = make([]*models.Participant, 0)

	//get all record
	var cur *mongo.Cursor
	if condition == nil {
		result, err := collection.Find(ctx, bson.D{})
		if err != nil {
			return nil, err
		}
		defer result.Close(ctx)
		cur = result
	} else {
		result, err := collection.Find(ctx, condition)
		if err != nil {
			return nil, err
		}
		defer result.Close(ctx)
		cur = result
	}

	//map data to target variable
	for cur.Next(ctx) {
		var participant models.Participant
		cur.Decode(&participant)
		participants = append(participants, &participant)
	}
	//response data to client
	if participants == nil {
		return make([]*models.Participant, 0), nil
	}
	return participants, nil
}

/*FindOne: get one record from a collection  */
func (u *ParticipantRepository) FindOne(filter bson.M) (*models.Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("participant")
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

	participant := models.Participant{}
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
func (u *ParticipantRepository) Create(newParticipant models.Participant) (*models.Participant, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionParticipantName)
	defer cancel()

	//convert to bson.M
	currentTime := time.Now()
	participant := models.Participant{
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		IsValid:              false,
		IsAttended:           false,
		Event:                newParticipant.Event,
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

	return &models.Participant{
		ID:                   insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		IsValid:              false,
		IsAttended:           false,
		Event:                newParticipant.Event,
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
func (u ParticipantRepository) UpdateOne(filter bson.M, update bson.M) (*models.Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionParticipantName)
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
	participant, errQuery := u.FindOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return participant, nil
}

//DeleteOne func is to update one record from a collection
func (u ParticipantRepository) DeleteOne(filter bson.M) (*models.Participant, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionParticipantName)
	defer cancel()

	participant, errFind := u.FindOne(filter)
	if errFind != nil {
		return nil, errFind
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
