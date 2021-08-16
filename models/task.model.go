package models

import (
	"errors"
	"time"

	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Event     Event              `bson:"event" json:"event"`
	Name      string             `bson:"name" json:"name"`
	User      User               `bson:"user" json:"user"`
	Type      string             `bson:"type" json:"type"`
	StartDate time.Time          `bson:"startDate" json:"startDate"`
	EndDate   time.Time          `bson:"endDate" json:"endDate"`
}

/* Input Type */
type NewTask struct {
	Event     string    `bson:"event" json:"event"`
	Name      string    `bson:"name" json:"name"`
	User      string    `bson:"user" json:"user"`
	Type      string    `bson:"type" json:"type"`
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
}

/* Model Function */
/* GetAll: get all data based on condition*/
func (u *Task) GetAll() ([]*Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("tasks")
	defer cancel()

	//create an empty array to store all fields from collection
	var tasks []*Task = make([]*Task, 0)

	//get all record
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	//map data to target variable
	for cur.Next(ctx) {
		var task Task
		cur.Decode(&task)
		tasks = append(tasks, &task)
	}
	//response data to client
	if tasks == nil {
		return make([]*Task, 0), nil
	}
	return tasks, nil
}

/*GetOne: get one record from a collection  */
func (u *Task) GetOne(filter bson.M) (*Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("task")
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

	task := Task{}
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
func (u *Task) Create(newTask NewTask) (*Task, error) {

	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("tasks")
	defer cancel()

	//get event, user
	var event *Event
	var user *User

	if targetUser, err := UserInstance.GetOne(bson.M{"_id": newTask.User}); err != nil {
		return nil, err
	} else if targetUser == nil {
		return nil, errors.New("invalid User id")
	} else {
		user = targetUser
	}
	if targetEvent, err := EventInstance.GetOne(bson.M{"_id": newTask.Event}); err != nil {
		return nil, err
	} else if targetEvent == nil {
		return nil, errors.New("invalid event id")
	} else {
		event = targetEvent
	}

	//convert to bson.M
	currentTime := time.Now()
	task := Task{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Event:     *event,
		Name:      newTask.Name,
		User:      *user,
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

	return &Task{
		ID:        insertResult.InsertedID.(primitive.ObjectID),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Event:     *event,
		Name:      newTask.Name,
		User:      *user,
		Type:      newTask.Type,
		StartDate: newTask.StartDate,
		EndDate:   newTask.EndDate,
	}, nil
}

/*UpdateOne: update one record from a collection*/
func (u Task) UpdateOne(filter bson.M, update bson.M) (*Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("tasks")
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
	task, errQuery := u.GetOne(filter)
	if errQuery != nil {
		return nil, errQuery
	}

	return task, nil
}

//DeleteOne func is to update one record from a collection
func (u Task) DeleteOne(filter bson.M) (*Task, error) {
	//get a collection , context, cancel func
	collection, ctx, cancel := createContextAndTargetCol("tasks")
	defer cancel()

	task, errGet := u.GetOne(filter)
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

	return task, nil
}

//export an instance Task model
var TaskInstance = Task{}
