package services

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TaskServiceName = "TaskServiceName"

type TaskService struct {
	TaskRepository *TaskRepository
}

/* GetAll: get all data based on condition*/
func (u *TaskService) GetAll(condition bson.M) ([]*models.Task, error) {
	return u.TaskRepository.FindAll(condition)
}

/*GetOne: get one record from a collection  */
func (u *TaskService) GetOne(filter bson.M) (*models.Task, error) {
	return u.TaskRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *TaskService) Create(newTask model.NewTask) (*models.Task, error) {

	//get event, user
	eventId, err := primitive.ObjectIDFromHex(*newTask.EventID)
	if err != nil {
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(newTask.UserID)
	if err != nil {
		return nil, err
	}

	//convert to bson.M
	currentTime := time.Now()
	task := models.Task{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Event:     eventId,
		Name:      newTask.Name,
		User:      userId,
		Type:      newTask.Type,
		StartDate: newTask.StartDate,
		EndDate:   newTask.EndDate,
	}
	return u.TaskRepository.Create(&task)

}

/*UpdateOne: update one record from a collection*/
func (u TaskService) UpdateOne(filter bson.M, update model.UpdateTask) (*models.Task, error) {
	bsonUpdate, err := utilities.InterfaceToBsonM(update)
	if err != nil {
		return nil, err
	}
	return u.TaskRepository.UpdateOne(filter, bsonUpdate)
}

//DeleteOne func is to update one record from a collection
func (u TaskService) DeleteOne(filter bson.M) (*models.Task, error) {
	return u.TaskRepository.DeleteOne(filter)
}

//validation
func (u *TaskService) ValidateNewTask(newTask model.NewTask) error {
	return validation.ValidateStruct(&newTask,
		validation.Field(&newTask.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			event, err := u.GetOne(bson.M{"name": name.(string)})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			if event != nil {
				return errors.New("name already existed")
			}
			return nil

		})),
		validation.Field(&newTask.UserID, validation.Required.Error("user id must not be blanked")),
		validation.Field(&newTask.Type, validation.Required.Error("type must not be blanked")),
		validation.Field(&newTask.StartDate, validation.Required.Error("start date must not be blanked")),
		validation.Field(&newTask.EndDate, validation.Required.Error("end date must not be blanked")),
	)
}

func (u *TaskService) ValidateUpdateTask(id string, updateTask model.UpdateTask) error {
	return validation.ValidateStruct(&updateTask,
		validation.Field(&updateTask.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			//convert string id to object id
			objectId, err := utilities.ConvertStringIdToObjectID(id)
			if err != nil {
				return err
			}
			//get current task
			currentTask, err := u.GetOne(bson.M{"_id": objectId})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			//check email existed or not
			if task, err := u.GetOne(bson.M{"name": name.(string)}); err != nil {
				if _, ok := err.(*helpers.ErrNotFound); ok {
					return nil
				} else {
					return err
				}
			} else {
				if task.Name != currentTask.Name {
					return errors.New("name already existed")
				} else {
					return nil
				}
			}

		})),
		validation.Field(&updateTask.UserID, validation.Required.Error("user id must not be blanked")),
		validation.Field(&updateTask.Type, validation.Required.Error("type must not be blanked")),
		validation.Field(&updateTask.StartDate, validation.Required.Error("start date must not be blanked")),
		validation.Field(&updateTask.EndDate, validation.Required.Error("end date must not be blanked")),
	)
}
