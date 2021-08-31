package services

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

var EventTypeServiceName = "EventTypeServiceName"

type EventTypeService struct {
	EventTypeRepository *EventTypeRepository
}

/* GetAll: get all data based on condition*/
func (u *EventTypeService) GetAll(condition bson.M) ([]*models.EventType, error) {
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
func (u EventTypeService) UpdateOne(filter bson.M, update model.UpdateEventType) (*models.EventType, error) {
	bsonUpdate, err := utilities.InterfaceToBsonM(update)
	if err != nil {
		return nil, err
	}
	return u.EventTypeRepository.UpdateOne(filter, bsonUpdate)
}

//DeleteOne func is to update one record from a collection
func (u EventTypeService) DeleteOne(filter bson.M) (*models.EventType, error) {
	return u.EventTypeRepository.DeleteOne(filter)
}

//validation
func (u *EventTypeService) ValidateNewEventType(newEventType model.NewEventType) error {
	return validation.ValidateStruct(&newEventType,
		validation.Field(&newEventType.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			event, err := u.GetOne(bson.M{"name": name.(string)})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			if event != nil {
				return errors.New("name already existed")
			}
			return nil

		})),
	)
}

func (u *EventTypeService) ValidateUpdateEventType(id string, updateEventType model.UpdateEventType) error {
	return validation.ValidateStruct(&updateEventType,
		validation.Field(&updateEventType.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			//convert string id to object id
			objectId, err := utilities.ConvertStringIdToObjectID(id)
			if err != nil {
				return err
			}
			//get current eventType
			currentEventType, err := u.GetOne(bson.M{"_id": objectId})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			//check email existed or not
			if eventType, err := u.GetOne(bson.M{"name": name.(string)}); err != nil {
				if _, ok := err.(*helpers.ErrNotFound); ok {
					return nil
				} else {
					return err
				}
			} else {
				if eventType.Name != currentEventType.Name {
					return errors.New("name already existed")
				} else {
					return nil
				}
			}

		})),
	)
}
