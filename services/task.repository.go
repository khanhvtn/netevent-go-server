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

var TaskRepositoryName = "TaskRepositoryName"

type TaskRepository struct {
	MongoCN *database.MongoInstance
}

/* createContextAndTargetCol: create and return targeted collection based on collection name */
func (u *TaskRepository) createContextAndTargetCol(colName string) (col *mongo.Collection,
	ctx context.Context,
	cancel context.CancelFunc) {
	col = u.MongoCN.Db.Collection(colName)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	return
}

/* FindAll: get all data based on condition*/
func (u *TaskRepository) FindAll(condition *bson.M) ([]*models.Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionTaskName)
	defer cancel()

	//create an empty array to store all fields from collection
	var tasks []*models.Task = make([]*models.Task, 0)

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
		var task models.Task
		cur.Decode(&task)
		tasks = append(tasks, &task)
	}
	//response data to client
	if tasks == nil {
		return make([]*models.Task, 0), nil
	}
	return tasks, nil
}

/*FindOne: get one record from a collection  */
func (u *TaskRepository) FindOne(filter bson.M) (*models.Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol("task")
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

	task := models.Task{}
	//Decode record into result
	if err := collection.FindOne(ctx, filter).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			//return nil data when id is not existed.
			return nil, nil
		}
		//return err if there is a system error
		return nil, err
	}

	return &task, nil
}

/*Create: create a new record to a collection*/
func (u *TaskRepository) Create(newTask *models.Task) (*models.Task, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionTaskName)
	defer cancel()

	//convert to bson.M
	currentTime := time.Now()
	task := models.Task{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Event:     newTask.Event,
		Name:      newTask.Name,
		User:      newTask.User,
		Type:      newTask.Type,
		StartDate: newTask.StartDate,
		EndDate:   newTask.EndDate,
	}
	newData, err := utilities.InterfaceToBsonM(task)
	if err != nil {
		return nil, err
	}

	//create user in database
	insertResult, err := collection.InsertOne(ctx, newData)
	if err != nil {
		return nil, err
	}

	return &models.Task{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Event:     newTask.Event,
		Name:      newTask.Name,
		User:      newTask.User,
		Type:      newTask.Type,
		StartDate: newTask.StartDate,
		EndDate:   newTask.EndDate,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u TaskRepository) UpdateOne(filter bson.M, update bson.M) (*models.Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionTaskName)
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
	task, errQuery := u.FindOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return task, nil
}

//DeleteOne func is to update one record from a collection
func (u TaskRepository) DeleteOne(filter bson.M) (*models.Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := u.createContextAndTargetCol(models.CollectionTaskName)
	defer cancel()

	task, errFind := u.FindOne(filter)
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

	return task, nil
}
