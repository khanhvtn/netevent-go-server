package services

import (
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

var EventTypeServiceName = "EventTypeServiceName"

type EventTypeService struct {
	EventTypeRepository *EventTypeRepository
}

/* GetAll: get all data based on condition*/
func (u *EventTypeService) GetAll(condition *bson.M) ([]*models.EventType, error) {
	return u.EventTypeRepository.Find(condition)
}

/*GetOne: get one record from a collection  */
func (u *EventTypeService) GetOne(filter bson.M) (*models.EventType, error) {
	return u.EventTypeRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *EventTypeService) Create(newEventType model.NewEventType) (*models.EventType, error) {
	return u.EventTypeRepository.Create(newEventType)
}

/*UpdateOne: update one record from a collection*/
func (u EventTypeService) UpdateOne(filter bson.M, update bson.M) (*models.EventType, error) {
	return u.EventTypeRepository.UpdateOne(filter, update)
}

//DeleteOne func is to update one record from a collection
func (u EventTypeService) DeleteOne(filter bson.M) (*models.EventType, error) {
	return u.EventTypeRepository.DeleteOne(filter)
}
