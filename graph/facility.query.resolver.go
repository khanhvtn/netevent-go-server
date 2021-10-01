package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Facilities(ctx context.Context, filter model.FacilityFilter) ([]*model.Facility, error) {
	service := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	facilities, err := service.GetAll(filter)
	if err != nil {
		return nil, err
	}
	results := make([]*model.Facility, 0)
	for _, facility := range facilities {
		mappedFacility, err := r.mapFacility(facility)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedFacility)
	}
	return results, nil
}

func (r *queryResolver) Facility(ctx context.Context, id string) (*model.Facility, error) {
	service := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	//get facility based specific id
	facility, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapFacility(facility)
	if err != nil {
		return nil, err
	}
	return result, nil
}
