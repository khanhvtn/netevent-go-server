package services

import (
	"errors"
	"sync"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
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
func (u *EventService) GetAll(condition bson.M) ([]*models.Event, error) {
	return u.EventRepository.FindAll(condition)
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
			if err := u.rollbackForCreateEvent(taskIds, models.CollectionTaskName); err != nil {
				return nil, err
			}
			if err := u.rollbackForCreateEvent(facilityHistoryIds, models.CollectionFacilityHistoryName); err != nil {
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
		if err := u.rollbackForCreateEvent(taskIds, models.CollectionTaskName); err != nil {
			return nil, err
		}
		if err := u.rollbackForCreateEvent(facilityHistoryIds, models.CollectionFacilityHistoryName); err != nil {
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
func (u EventService) UpdateOne(filter bson.M, update model.UpdateEvent) (*models.Event, error) {
	/*
		Create backup for event, tasks, and facility histories.
		Fetch current tasks and facilityHistories correspond with event id
		Compare current tasks and facilityHistories with newUpdate:
		-if the number of ids of current tasks and facilityHistories are similar to newUpdate:
			-update them with correspond ids and return ids to chan
		-if there are new tasks or facilityHistories:
			-create them and return the id to chan
		-if there are any errors:
			-return errors to error chan.
	*/

	taskIds := make([]primitive.ObjectID, 0)
	facilityHistoryIds := make([]primitive.ObjectID, 0)
	lenObjectIdsChan := len(update.FacilityHistories) + len(update.Tasks)
	objectIdsChan := make(chan map[string]primitive.ObjectID, lenObjectIdsChan)
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	//get current tasks and facilityHistories correspond with event id
	currentEvent, err := u.GetOne(filter)
	if err != nil {
		return nil, err
	}
	if currentEvent == nil {
		return nil, errors.New("event id not found")
	}
	backupTasks, err := u.TaskRepository.FindAll(bson.M{"event": currentEvent.ID})
	if err != nil {
		return nil, err
	}
	backupFacilityHistories, err := u.FacilityHistoryRepository.FindAll(bson.M{"event": currentEvent.ID})
	if err != nil {
		return nil, err
	}

	//update task ids
	wg.Add(1)
	go func(update model.UpdateEvent, currentEvent *models.Event, objectIdsChan chan<- map[string]primitive.ObjectID, errChan chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()
		currentTaskIds := currentEvent.Tasks
		var errDB error = nil
		for _, newTask := range update.Tasks {
			userID, err := primitive.ObjectIDFromHex(newTask.UserID)
			if err != nil {
				errDB = err
				break
			}
			if newTask.ID == nil {
				createdTask, err := u.TaskRepository.Create(&models.Task{
					Event:     currentEvent.ID,
					Name:      newTask.Name,
					User:      userID,
					Type:      newTask.Type,
					StartDate: newTask.StartDate,
					EndDate:   newTask.EndDate,
				})
				if err != nil {
					errDB = err
					break
				}
				objectIdsChan <- map[string]primitive.ObjectID{models.CollectionTaskName: createdTask.ID}
			} else if primitive.IsValidObjectID(*newTask.ID) {
				taskID, err := primitive.ObjectIDFromHex(*newTask.ID)
				if err != nil {
					errDB = err
					break
				}

				targetTask, err := u.TaskRepository.FindOne(bson.M{"_id": taskID})
				if err != nil {
					errDB = err
					break
				}

				targetTask.Event = currentEvent.ID
				targetTask.Name = newTask.Name
				targetTask.User = userID
				targetTask.Type = newTask.Type
				targetTask.StartDate = newTask.StartDate
				targetTask.EndDate = newTask.EndDate
				targetTask.UpdatedAt = time.Now()
				bsonTask, err := utilities.InterfaceToBsonM(targetTask)
				if err != nil {
					errDB = err
					break
				}

				updatedTask, err := u.TaskRepository.UpdateOne(bson.M{"_id": taskID}, bsonTask)
				if err != nil {
					errDB = err
					break
				}
				objectIdsChan <- map[string]primitive.ObjectID{models.CollectionTaskName: updatedTask.ID}
				currentTaskIds = u.removeObjectId(currentTaskIds, func(oi primitive.ObjectID) bool { return oi != updatedTask.Event })
			} else {
				errDB = errors.New("invalid task id")
				break
			}
		}
		//remove all task ids in current event if they are not in updated request.
		if errDB != nil {
			for _, v := range currentTaskIds {
				_, err := u.TaskRepository.DeleteOne(bson.M{"_id": v})
				if err != nil {
					errDB = err
					break
				}
			}
		}
		errChan <- errDB

	}(update, currentEvent, objectIdsChan, errChan, &wg)

	//update facility history ids
	wg.Add(1)
	go func(update model.UpdateEvent, currentEvent *models.Event, objectIdsChan chan<- map[string]primitive.ObjectID, errChan chan<- error, wg *sync.WaitGroup) {
		currentFacilityHistoryIds := currentEvent.FacilityHistories
		defer wg.Done()
		var errDB error = nil
		for _, newFacilityHistory := range update.FacilityHistories {
			facilityID, err := primitive.ObjectIDFromHex(newFacilityHistory.FacilityID)
			if err != nil {
				errDB = err
				break
			}

			if newFacilityHistory.ID == nil {
				createdFacilityHistory, err := u.FacilityHistoryRepository.Create(&models.FacilityHistory{
					Event:      currentEvent.ID,
					Facility:   facilityID,
					BorrowDate: newFacilityHistory.BorrowDate,
					ReturnDate: newFacilityHistory.ReturnDate,
				})
				if err != nil {
					errDB = err
					break
				}
				objectIdsChan <- map[string]primitive.ObjectID{models.CollectionFacilityHistoryName: createdFacilityHistory.ID}
			} else if primitive.IsValidObjectID(*newFacilityHistory.ID) {
				facilityHistoryID, err := primitive.ObjectIDFromHex(*newFacilityHistory.ID)
				if err != nil {
					errDB = err
					break
				}
				targetFacilityHistory, err := u.FacilityHistoryRepository.FindOne(bson.M{"_id": facilityHistoryID})
				if err != nil {
					errDB = err
					break
				}
				targetFacilityHistory.Event = currentEvent.ID
				targetFacilityHistory.Facility = facilityID
				targetFacilityHistory.BorrowDate = newFacilityHistory.BorrowDate
				targetFacilityHistory.ReturnDate = newFacilityHistory.ReturnDate
				targetFacilityHistory.UpdatedAt = time.Now()
				bsonFacilityHistory, err := utilities.InterfaceToBsonM(targetFacilityHistory)
				if err != nil {
					errDB = err
					break
				}
				updatedFacilityHistory, err := u.FacilityHistoryRepository.UpdateOne(bson.M{"_id": facilityHistoryID}, bsonFacilityHistory)
				if err != nil {
					errDB = err
					break
				}
				objectIdsChan <- map[string]primitive.ObjectID{models.CollectionFacilityHistoryName: updatedFacilityHistory.ID}
				currentFacilityHistoryIds = u.removeObjectId(currentFacilityHistoryIds, func(oi primitive.ObjectID) bool { return oi != updatedFacilityHistory.Event })
			} else {
				errDB = errors.New("invalid facility history id")
				break
			}
		}
		//remove all facility histories ids in current event if they are not in updated request.
		if errDB != nil {
			for _, v := range currentFacilityHistoryIds {
				_, err := u.FacilityHistoryRepository.DeleteOne(bson.M{"_id": v})
				if err != nil {
					errDB = err
					break
				}
			}
		}
		errChan <- errDB
	}(update, currentEvent, objectIdsChan, errChan, &wg)
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
			if err := u.rollbackForUpdateEvent(backupTasks, backupFacilityHistories, taskIds, facilityHistoryIds); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	evenTypeID, err := primitive.ObjectIDFromHex(update.EventTypeID)
	if err != nil {
		return nil, err
	}
	ownerID, err := primitive.ObjectIDFromHex(update.OwnerID)
	if err != nil {
		return nil, err
	}
	var reviewerID *primitive.ObjectID = nil
	if update.Reviewer != nil {
		objectId, err := primitive.ObjectIDFromHex(*update.Reviewer)
		if err != nil {
			return nil, err
		}
		reviewerID = &objectId
	}

	//convert to bson.M
	currentTime := time.Now()
	event := models.Event{
		Tags:                  update.Tags,
		IsApproved:            update.IsApproved,
		Reviewer:              reviewerID,
		IsFinished:            update.IsFinished,
		Tasks:                 taskIds,
		FacilityHistories:     facilityHistoryIds,
		Name:                  update.Name,
		Language:              update.Language,
		EventType:             evenTypeID,
		Mode:                  update.Mode,
		Location:              update.Location,
		Accommodation:         update.Accommodation,
		RegistrationCloseDate: update.RegistrationCloseDate,
		StartDate:             update.StartDate,
		EndDate:               update.EndDate,
		MaxParticipants:       update.MaxParticipants,
		Description:           update.Description,
		Owner:                 ownerID,
		Budget:                update.Budget,
		Image:                 update.Image,
		IsDeleted:             update.IsDeleted,
		CreatedAt:             currentEvent.CreatedAt,
		UpdatedAt:             currentTime,
	}
	bsonEvent, err := utilities.InterfaceToBsonM(event)
	if err != nil {
		return nil, err
	}

	updatedEvent, err := u.EventRepository.UpdateOne(filter, bsonEvent)
	if err != nil {
		if err := u.rollbackForUpdateEvent(backupTasks, backupFacilityHistories, taskIds, facilityHistoryIds); err != nil {
			return nil, err
		}
		return nil, err
	}
	return updatedEvent, nil
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
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
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
			//convert string id to object id
			objectId, err := utilities.ConvertStringIdToObjectID(id)
			if err != nil {
				return err
			}
			//get current event
			currentEvent, err := u.GetOne(bson.M{"_id": objectId})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			//check email existed or not
			if event, err := u.GetOne(bson.M{"name": name.(string)}); err != nil {
				if _, ok := err.(*helpers.ErrNotFound); ok {
					return nil
				} else {
					return err
				}
			} else {
				if event.Name != currentEvent.Name {
					return errors.New("name already existed")
				} else {
					return nil
				}
			}

		})),
		validation.Field(&updateEvent.Tags, validation.Required.Error("tags must not be blanked")),
		validation.Field(&updateEvent.Tasks, validation.Required.Error("tasks must not be blanked")),
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

/*rollbackForCreateEvent: remove all the tasks and facilityHistories that created to put into event when the event create fail*/
func (u *EventService) rollbackForCreateEvent(objectIds []primitive.ObjectID, collectionName string) error {
	if collectionName == models.CollectionTaskName {
		for _, v := range objectIds {
			_, err := u.TaskRepository.DeleteOne(bson.M{"_id": v})
			if err != nil {
				return err
			}
		}
	} else {
		for _, v := range objectIds {
			_, err := u.FacilityHistoryRepository.DeleteOne(bson.M{"_id": v})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*rollbackForUpdateEvent: rollback to previous state when the event update fail*/
func (u *EventService) rollbackForUpdateEvent(backupTasks []*models.Task, backupFacilityHistories []*models.FacilityHistory, createdTasks []primitive.ObjectID, createdFacilityHistories []primitive.ObjectID) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	if err := u.rollbackForCreateEvent(createdTasks, models.CollectionTaskName); err != nil {
		return err
	}
	if err := u.rollbackForCreateEvent(createdTasks, models.CollectionTaskName); err != nil {
		return err
	}
	wg.Add(1)
	go func(errChan chan<- error, backupTasks []*models.Task, wg *sync.WaitGroup) {
		defer wg.Done()
		var errDB error = nil
		for _, v := range backupTasks {
			_, err := u.TaskRepository.Create(v)
			if err != nil {
				errDB = err
				break
			}
		}
		errChan <- errDB
	}(errChan, backupTasks, &wg)
	wg.Add(1)
	go func(errChan chan<- error, backupFacilityHistories []*models.FacilityHistory, wg *sync.WaitGroup) {
		defer wg.Done()
		var errDB error = nil
		for _, v := range backupFacilityHistories {
			_, err := u.FacilityHistoryRepository.Create(v)
			if err != nil {
				errDB = err
				break
			}
		}
		errChan <- errDB
	}(errChan, backupFacilityHistories, &wg)
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

/*createTasksForEvent: create multiple tasks when creating an event  */
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

/*createTasksForEvent: create multiple facility histories when creating an event  */
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

/*updateEventId: update event id to target collection name with correspond id*/
func (u *EventService) updateEventId(listTargetIds []primitive.ObjectID, eventId primitive.ObjectID, collectionName string) error {
	if collectionName == models.CollectionTaskName {
		for _, v := range listTargetIds {
			_, err := u.TaskRepository.UpdateOne(bson.M{"_id": v}, bson.M{"event": eventId})
			if err != nil {
				return err
			}
		}
	} else {
		for _, v := range listTargetIds {
			_, err := u.FacilityHistoryRepository.UpdateOne(bson.M{"_id": v}, bson.M{"event": eventId})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*removeObjectId:  removespecific objectID in a slice objectIDs and return new slice */
func (u *EventService) removeObjectId(vs []primitive.ObjectID, f func(primitive.ObjectID) bool) []primitive.ObjectID {
	tmpSlice := make([]primitive.ObjectID, 0)
	for _, v := range vs {
		if f(v) {
			tmpSlice = append(tmpSlice, v)
		}
	}
	return tmpSlice
}
