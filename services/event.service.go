package services

import (
	"errors"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var EventServiceName = "EventServiceName"

type EventService struct {
	EventRepository           *EventRepository
	TaskRepository            *TaskRepository
	FacilityHistoryRepository *FacilityHistoryRepository
}

/* GetAll: get all data based on condition*/
func (u *EventService) GetAll(condition *bson.M) ([]*models.Event, error) {
	return u.EventRepository.FindAll(nil)
}

/*GetOne: get one record from a collection  */
func (u *EventService) GetOne(filter bson.M) (*models.Event, error) {
	return u.EventRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *EventService) Create(newEvent model.NewEvent) (*models.Event, error) {

	//create and collect ids for task and facility history
	taskIds := make([]primitive.ObjectID, 0)
	facilityHistoryIds := make([]primitive.ObjectID, 0)
	lenObjectIdsChan := len(newEvent.FacilityHistories) + len(newEvent.Tasks)
	objectIdsChan := make(chan map[string]primitive.ObjectID, lenObjectIdsChan)
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go u.createTasksForEvent(newEvent.Tasks, objectIdsChan, errChan, &wg)
	wg.Add(1)
	go u.createFacilityHistoriesForEvent(newEvent.FacilityHistories, objectIdsChan, errChan, &wg)
	wg.Wait()
	close(objectIdsChan)
	close(errChan)

	for mapObjectId := range objectIdsChan {

		if value, ok := mapObjectId[models.CollectionTaskName]; ok {
			taskIds = append(taskIds, value)
		} else {
			value = mapObjectId[models.CollectionFacilityHistoryName]
			facilityHistoryIds = append(facilityHistoryIds, value)
		}
	}
	for err := range errChan {
		if err != nil {
			if err := u.cleanUpObjectIds(taskIds, models.CollectionTaskName); err != nil {
				return nil, err
			}
			if err := u.cleanUpObjectIds(facilityHistoryIds, models.CollectionFacilityHistoryName); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	evenTypeID, err := primitive.ObjectIDFromHex(newEvent.EventTypeID)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(newEvent.OwnerID)
	if err != nil {
		return nil, err
	}

	//convert to bson.M
	currentTime := time.Now()
	event := models.Event{
		Tags:                  newEvent.Tags,
		IsApproved:            false,
		Reviewer:              nil,
		IsFinished:            false,
		Tasks:                 taskIds,
		FacilityHistories:     facilityHistoryIds,
		Name:                  newEvent.Name,
		Language:              newEvent.Language,
		EventType:             evenTypeID,
		Mode:                  newEvent.Mode,
		Location:              newEvent.Location,
		Accommodation:         newEvent.Accommodation,
		RegistrationCloseDate: newEvent.RegistrationCloseDate,
		StartDate:             newEvent.StartDate,
		EndDate:               newEvent.EndDate,
		MaxParticipants:       newEvent.MaxParticipants,
		Description:           newEvent.Description,
		Owner:                 ownerID,
		Budget:                newEvent.Budget,
		Image:                 newEvent.Image,
		IsDeleted:             false,
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
	}
	//create event
	createdEvent, err := u.EventRepository.Create(event)
	if err != nil {
		if err := u.cleanUpObjectIds(taskIds, models.CollectionTaskName); err != nil {
			return nil, err
		}
		if err := u.cleanUpObjectIds(facilityHistoryIds, models.CollectionFacilityHistoryName); err != nil {
			return nil, err
		}
		return nil, err
	}
	//update event id back to facilityhistories and tasks
	u.updateEventId(taskIds, createdEvent.ID, "tasks")
	u.updateEventId(facilityHistoryIds, createdEvent.ID, "facilityHistories")

	return createdEvent, nil
}

/*UpdateOne: update one record from a collection*/
func (u EventService) UpdateOne(filter bson.M, update bson.M) (*models.Event, error) {
	return u.EventRepository.UpdateOne(filter, update)
}

//DeleteOne func is to update one record from a collection
func (u EventService) DeleteOne(filter bson.M) (*models.Event, error) {
	return u.EventRepository.DeleteOne(filter)
}

//validation
func (u *EventService) ValidateNewEvent(newEvent model.NewEvent) error {
	return validation.ValidateStruct(&newEvent,
		validation.Field(&newEvent.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			event, err := u.GetOne(bson.M{"name": name.(string)})
			if err != nil {
				return err
			}
			if event != nil {
				return errors.New("name already existed")
			}
			return nil

		})),
		validation.Field(&newEvent.Tags, validation.Required.Error("tags must not be blanked")),
		validation.Field(&newEvent.Tasks, validation.Required.Error("tasks password must not be blanked")),
		validation.Field(&newEvent.FacilityHistories, validation.Required.Error("facility histories must not be blanked")),
		validation.Field(&newEvent.Language, validation.Required.Error("language must not be blanked")),
		validation.Field(&newEvent.EventTypeID, validation.Required.Error("event type must not be blanked")),
		validation.Field(&newEvent.Mode, validation.Required.Error("mode must not be blanked")),
		validation.Field(&newEvent.Location, validation.Required.Error("location must not be blanked")),
		validation.Field(&newEvent.Accommodation, validation.Required.Error("accommodation must not be blanked")),
		validation.Field(&newEvent.StartDate, validation.Required.Error("start date must not be blanked")),
		validation.Field(&newEvent.EndDate, validation.Required.Error("end date must not be blanked")),
		validation.Field(&newEvent.RegistrationCloseDate, validation.Required.Error("registration close date must not be blanked")),
		validation.Field(&newEvent.MaxParticipants, validation.Required.Error("max participants must not be blanked")),
		validation.Field(&newEvent.Description, validation.Required.Error("description must not be blanked")),
		validation.Field(&newEvent.OwnerID, validation.Required.Error("owner must not be blanked")),
		validation.Field(&newEvent.Budget, validation.Required.Error("budget must not be blanked")),
		validation.Field(&newEvent.Image, validation.Required.Error("image must not be blanked")),
	)
}

func (u *EventService) ValidateUpdateEvent(id string, updateEvent model.UpdateEvent) error {
	return validation.ValidateStruct(&updateEvent,
		validation.Field(&updateEvent.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			//get current event
			currentEvent, err := u.GetOne(bson.M{"_id": id})
			if err != nil {
				return err
			}
			//check email existed or not
			event, err := u.GetOne(bson.M{"name": name.(string)})
			if err != nil {
				return err
			}
			if event != nil && event.Name != currentEvent.Name {
				return errors.New("name already existed")
			}
			return nil

		})),
		validation.Field(&updateEvent.Tags, validation.Required.Error("tags must not be blanked")),
		validation.Field(&updateEvent.Tasks, validation.Required.Error("tasks password must not be blanked")),
		validation.Field(&updateEvent.FacilityHistories, validation.Required.Error("facility histories must not be blanked")),
		validation.Field(&updateEvent.Language, validation.Required.Error("language must not be blanked")),
		validation.Field(&updateEvent.EventTypeID, validation.Required.Error("event type must not be blanked")),
		validation.Field(&updateEvent.Mode, validation.Required.Error("mode must not be blanked")),
		validation.Field(&updateEvent.Location, validation.Required.Error("location must not be blanked")),
		validation.Field(&updateEvent.Accommodation, validation.Required.Error("accommodation must not be blanked")),
		validation.Field(&updateEvent.StartDate, validation.Required.Error("start date must not be blanked")),
		validation.Field(&updateEvent.EndDate, validation.Required.Error("end date must not be blanked")),
		validation.Field(&updateEvent.RegistrationCloseDate, validation.Required.Error("registration close date must not be blanked")),
		validation.Field(&updateEvent.MaxParticipants, validation.Required.Error("max participants must not be blanked")),
		validation.Field(&updateEvent.Description, validation.Required.Error("description must not be blanked")),
		validation.Field(&updateEvent.OwnerID, validation.Required.Error("owner must not be blanked")),
		validation.Field(&updateEvent.Budget, validation.Required.Error("budget must not be blanked")),
		validation.Field(&updateEvent.Image, validation.Required.Error("image must not be blanked")),
	)
}

//
func (u *EventService) cleanUpObjectIds(objectIds []primitive.ObjectID, collectionName string) error {
	if collectionName == models.CollectionTaskName {
		for _, v := range objectIds {
			_, err := u.TaskRepository.DeleteOne(bson.M{"_id": v.Hex()})
			if err != nil {
				return err
			}
		}
	} else {
		for _, v := range objectIds {
			_, err := u.FacilityHistoryRepository.DeleteOne(bson.M{"_id": v.Hex()})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (u *EventService) createTasksForEvent(tasks []*model.NewTask, objectIdsChan chan<- map[string]primitive.ObjectID, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	var errDB error = nil
	for _, task := range tasks {
		userId, err := primitive.ObjectIDFromHex(task.UserID)
		if err != nil {
			errDB = err
			break
		}

		newTask, err := u.TaskRepository.Create(&models.Task{
			Name:      task.Name,
			Type:      task.Type,
			User:      userId,
			StartDate: task.StartDate,
			EndDate:   task.EndDate,
		})
		if err != nil {
			errDB = err
			break
		}
		objectIdsChan <- map[string]primitive.ObjectID{models.CollectionTaskName: newTask.ID}
	}
	errChan <- errDB
}
func (u *EventService) createFacilityHistoriesForEvent(facilityHistories []*model.NewFacilityHistory, objectIdsChan chan<- map[string]primitive.ObjectID, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	var errDB error = nil
	for _, facilityHistory := range facilityHistories {
		facilityId, err := primitive.ObjectIDFromHex(facilityHistory.FacilityID)
		if err != nil {
			errDB = err
			break
		}
		newFacilityHistory, err := u.FacilityHistoryRepository.Create(&models.FacilityHistory{
			Facility:   facilityId,
			BorrowDate: facilityHistory.BorrowDate,
			ReturnDate: facilityHistory.ReturnDate,
		})
		if err != nil {
			errDB = err
			break
		}
		objectIdsChan <- map[string]primitive.ObjectID{models.CollectionFacilityHistoryName: newFacilityHistory.ID}
	}
	errChan <- errDB
}

func (u *EventService) updateEventId(listTargetIds []primitive.ObjectID, eventId primitive.ObjectID, collectionName string) error {
	if collectionName == models.CollectionTaskName {
		for _, v := range listTargetIds {
			_, err := u.TaskRepository.UpdateOne(bson.M{"_id": v.Hex()}, bson.M{"event": eventId})
			if err != nil {
				return err
			}
		}
	} else {
		for _, v := range listTargetIds {
			_, err := u.FacilityHistoryRepository.UpdateOne(bson.M{"_id": v.Hex()}, bson.M{"event": eventId})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
