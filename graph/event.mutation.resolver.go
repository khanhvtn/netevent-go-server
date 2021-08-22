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
		return nil, err
	}
	newEvent, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, _ := mapEvent(newEvent)
	return results, nil
}
func (r *mutationResolver) UpdateEvent(ctx context.Context, id string, input model.UpdateEvent) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	//check input
	if err := service.ValidateUpdateEvent(id, input); err != nil {
		return nil, err
	}
	//cast UpdateEvent to bson.M type
	newUpdate, err := utilities.InterfaceToBsonM(input)
	if err != nil {
		return nil, err
	}
	updatedEvent, err := service.UpdateOne(bson.M{"_id": id}, newUpdate)
	if err != nil {
		return nil, err
	}
	results, _ := mapEvent(updatedEvent)
	return results, nil
}
func (r *mutationResolver) DeleteEvent(ctx context.Context, id string) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	deletedEvent, err := service.DeleteOne(bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	results, _ := mapEvent(deletedEvent)
	return results, nil
}
