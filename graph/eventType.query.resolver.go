package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) EventTypes(ctx context.Context, filter model.EventTypeFilter) (*model.EventTypeResponse, error) {
	service := r.di.Container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	eventTypes, pageInfo, err := service.GetAll(filter)
	if err != nil {
		return nil, err
	}
	results := make([]*model.EventType, 0)
	for _, eventType := range eventTypes {
		mappedEventType, err := r.mapEventType(eventType)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedEventType)
	}
	return &model.EventTypeResponse{PageInfo: pageInfo, EventTypes: results}, nil
}
func (r *queryResolver) EventType(ctx context.Context, id string) (*model.EventType, error) {
	service := r.di.Container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	//get event type based specific id
	eventType, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapEventType(eventType)
	if err != nil {
		return nil, err
	}
	return result, nil
}
