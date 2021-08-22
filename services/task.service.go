package services

import (
	"time"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TaskServiceName = "TaskServiceName"

type TaskService struct {
	TaskRepository *TaskRepository
}

/* GetAll: get all data based on condition*/
func (u *TaskService) GetAll(condition *bson.M) ([]*models.Task, error) {
	return u.TaskRepository.FindAll(condition)
}

/*GetOne: get one record from a collection  */
func (u *TaskService) GetOne(filter bson.M) (*models.Task, error) {
	return u.TaskRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *TaskService) Create(newTask *model.NewTask) (*models.Task, error) {

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
func (u TaskService) UpdateOne(filter bson.M, update bson.M) (*models.Task, error) {
	return u.TaskRepository.UpdateOne(filter, update)
}

//DeleteOne func is to update one record from a collection
func (u TaskService) DeleteOne(filter bson.M) (*models.Task, error) {
	return u.TaskRepository.DeleteOne(filter)
}
