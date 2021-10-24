package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	//check input
	if err := service.ValidateNewEvent(input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	newEvent, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapEvent(newEvent)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) UpdateEvent(ctx context.Context, id string, input model.UpdateEvent) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	//check input
	if err := service.ValidateUpdateEvent(id, input); err != nil {
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
	updatedEvent, err := service.UpdateOne(bson.M{"_id": objectId}, input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapEvent(updatedEvent)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) DeleteEvent(ctx context.Context, id string) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	deletedEvent, err := service.DeleteOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	results, err := r.mapEvent(deletedEvent)
	if err != nil {
		return nil, err
	}
	return results, nil
}
