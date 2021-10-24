package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateEventType(ctx context.Context, input model.NewEventType) (*model.EventType, error) {
	service := r.di.Container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	//check input
	if err := service.ValidateNewEventType(input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	newEventType, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapEventType(newEventType)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) UpdateEventType(ctx context.Context, id string, input model.UpdateEventType) (*model.EventType, error) {

	service := r.di.Container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	//check input
	if err := service.ValidateUpdateEventType(id, input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	updatedEventType, err := service.UpdateOne(bson.M{"_id": objectId}, input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapEventType(updatedEventType)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) DeleteEventType(ctx context.Context, id string) (*model.EventType, error) {
	service := r.di.Container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	deletedEventType, err := service.DeleteOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	results, err := r.mapEventType(deletedEventType)
	if err != nil {
		return nil, err
	}
	return results, nil
}
