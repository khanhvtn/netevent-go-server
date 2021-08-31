package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *queryResolver) Facilities(ctx context.Context) ([]*model.Facility, error) {
	service := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	users, err := service.GetAll(bson.M{})
	if err != nil {
		return nil, err
	}
	results := make([]*model.Facility, 0)
	for _, user := range users {
		mappedFacility, err := r.mapFacility(user)
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
	//get user based specific id
	user, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapFacility(user)
	if err != nil {
		return nil, err
	}
	return result, nil
}
