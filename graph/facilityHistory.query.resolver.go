package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *facilityHistoryResolver) Facility(ctx context.Context, obj *model.FacilityHistory) (*model.Facility, error) {
	facilityService := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	facilityHistoryService := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	facilityHistory, err := facilityHistoryService.GetOne(bson.M{"_id": obj.ID})
	if err != nil {
		return nil, err
	}
	facility, err := facilityService.GetOne(bson.M{"_id": facilityHistory.Facility})
	if err != nil {
		return nil, err
	}
	results, err := r.mapFacility(facility)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *facilityHistoryResolver) Event(ctx context.Context, obj *model.FacilityHistory) (*model.Event, error) {
	facilityHistoryService := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	eventService := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	facilityHistory, err := facilityHistoryService.GetOne(bson.M{"_id": obj.ID})
	if err != nil {
		return nil, err
	}
	event, err := eventService.GetOne(bson.M{"_id": facilityHistory.Event})
	if err != nil {
		return nil, err
	}
	results, err := r.mapEvent(event)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *queryResolver) FacilityHistories(ctx context.Context) ([]*model.FacilityHistory, error) {
	service := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	facilityHistories, err := service.GetAll(bson.M{})
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
func (r *queryResolver) FacilityHistory(ctx context.Context, id string) (*model.FacilityHistory, error) {
	service := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
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
	result, err := r.mapFacilityHistory(event)
	if err != nil {
		return nil, err
	}
	return result, nil
}
