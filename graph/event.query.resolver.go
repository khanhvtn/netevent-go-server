package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *eventResolver) Tasks(ctx context.Context, obj *model.Event) ([]*model.Task, error) {
	service := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	tasks, err := service.GetAll(bson.M{"event": obj.ID})
	if err != nil {
		return nil, err
	}
	results := make([]*model.Task, 0)
	for _, task := range tasks {
		mappedTask, err := r.mapTask(task)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedTask)
	}
	return results, nil
}

func (r *eventResolver) FacilityHistories(ctx context.Context, obj *model.Event) ([]*model.FacilityHistory, error) {
	service := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	facilityHistories, err := service.GetAll(bson.M{"event": obj.ID})

	if err != nil {
		return nil, err
	}
	results := make([]*model.FacilityHistory, 0)
	for _, facilityHistory := range facilityHistories {
		mappedFacilityHistory, err := r.mapFacilityHistory(facilityHistory)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedFacilityHistory)
	}

	return results, nil
}

func (r *eventResolver) EventType(ctx context.Context, obj *model.Event) (*model.EventType, error) {
	service := r.di.Container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	eventType, err := service.GetOne(bson.M{"_id": obj.EventType.ID})
	if err != nil {
		return nil, err
	}
	results, err := r.mapEventType(eventType)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *eventResolver) Owner(ctx context.Context, obj *model.Event) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	owner, err := service.GetOne(bson.M{"_id": obj.Owner.ID})
	if err != nil {
		return nil, err
	}
	results, err := r.mapUser(owner)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *queryResolver) Events(ctx context.Context) ([]*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	events, err := service.GetAll(bson.M{})
	if err != nil {
		return nil, err
	}
	results := make([]*model.Event, 0)
	for _, event := range events {
		mappedEvent, err := r.mapEvent(event)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedEvent)
	}
	return results, nil
}

func (r *queryResolver) Event(ctx context.Context, id string) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	//get event based specific id
	event, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapEvent(event)
	if err != nil {
		return nil, err
	}
	return result, nil
}
