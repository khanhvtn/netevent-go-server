package services

import (
	"errors"
	"math"
	"unsafe"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var EventTypeServiceName = "EventTypeServiceName"

type EventTypeService struct {
	EventTypeRepository *EventTypeRepository
}

/* GetAll: get all data based on condition*/
func (u *EventTypeService) GetAll(filter model.EventTypeFilter) ([]*models.EventType, *model.PageInfo, error) {
	//setup filter field
	finalFilter := bson.M{}
	filterByString := make([]bson.M, 0)
	opts := options.Find()
	//set filter name
	var keySearch string = ""
	if filter.DefaultFilter.Search != nil {
		keySearch = *filter.DefaultFilter.Search
	}
	filterByString = append(filterByString, bson.M{
		"name": bson.M{"$regex": primitive.Regex{Pattern: keySearch, Options: "i"}},
	})

	//for updatedAt
	if filter.DefaultFilter.UpdatedAtDateFrom != nil && filter.DefaultFilter.UpdatedAtDateTo != nil {
		finalFilter["updatedAt"] = bson.M{
			"$gte": filter.DefaultFilter.UpdatedAtDateFrom,
			"$lte": filter.DefaultFilter.UpdatedAtDateTo,
		}
	}
	//for createdAt
	if filter.DefaultFilter.CreatedAtDateFrom != nil && filter.DefaultFilter.CreatedAtDateTo != nil {
		finalFilter["createdAt"] = bson.M{
			"$gte": filter.DefaultFilter.CreatedAtDateFrom,
			"$lte": filter.DefaultFilter.CreatedAtDateTo,
		}
	}

	//set the number of record that will display
	var take int64 = 10 //take 10 records in default
	if filter.DefaultFilter.Take != nil {
		take = *(*int64)(unsafe.Pointer(filter.DefaultFilter.Take))
	}
	opts.SetLimit(take)

	//set isDeleted filter
	var isDeleted = false
	if filter.DefaultFilter.IsDeleted != nil {
		isDeleted = *filter.DefaultFilter.IsDeleted
	}
	finalFilter["isDeleted"] = isDeleted

	//set filter for string field
	finalFilter["$or"] = filterByString

	//get total number page
	errChan := make(chan error, 2)
	totalPageChan := make(chan int)
	var totalPage *int = nil

	go func(totalPageChan chan<- int, errChan chan<- error, finalFilter primitive.M, opts *options.FindOptions, take int64) {
		if eventTypes, err := u.EventTypeRepository.FindAll(finalFilter, opts); err != nil {
			errChan <- err
			close(totalPageChan)
		} else {
			totalEventTypes := len(eventTypes)
			totalPage := int(math.Ceil(float64(totalEventTypes) / float64(take)))
			errChan <- nil
			totalPageChan <- totalPage
		}

	}(totalPageChan, errChan, finalFilter, opts, take)

	//set paging
	var page int64 = 1 //target page 1 in default
	if filter.DefaultFilter.Page != nil {
		page = *(*int64)(unsafe.Pointer(filter.DefaultFilter.Page))
	}
	opts.SetSkip((page - 1) * take)

	//get eventTypes
	var eventTypes []*models.EventType = nil
	eventTypesChan := make(chan []*models.EventType)
	go func(eventTypesChan chan<- []*models.EventType, errChan chan<- error, finalFilter primitive.M, opts *options.FindOptions) {
		if eventTypes, err := u.EventTypeRepository.FindAll(finalFilter, opts); err != nil {
			errChan <- err
			close(eventTypesChan)
		} else {
			errChan <- nil
			eventTypesChan <- eventTypes
		}

	}(eventTypesChan, errChan, finalFilter, opts)

	if totalPageValue, ok := <-totalPageChan; ok {
		totalPage = &totalPageValue
	}
	if eventTypesValue, ok := <-eventTypesChan; ok {
		eventTypes = eventTypesValue
	}

	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, nil, err
		}
	}

	return eventTypes, &model.PageInfo{TotalPage: *totalPage, CurrentPage: int(page)}, nil
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
