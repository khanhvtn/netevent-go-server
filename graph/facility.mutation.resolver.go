package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateFacility(ctx context.Context, input model.NewFacility) (*model.Facility, error) {
	service := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	//check input
	if err := service.ValidateNewFacility(input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	newFacility, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapFacility(newFacility)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) UpdateFacility(ctx context.Context, id string, input model.UpdateFacility) (*model.Facility, error) {
	service := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	//check input
	if err := service.ValidateUpdateFacility(id, input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	//cast UpdateFacility to bson.M type
	newUpdate, err := utilities.InterfaceToBsonM(input)
	if err != nil {
		return nil, err
	}
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	updatedFacility, err := service.UpdateOne(bson.M{"_id": objectId}, newUpdate)
	if err != nil {
		return nil, err
	}
	results, err := r.mapFacility(updatedFacility)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) DeleteFacility(ctx context.Context, id string) (*model.Facility, error) {
	service := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	deletedFacility, err := service.DeleteOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	results, err := r.mapFacility(deletedFacility)
	if err != nil {
		return nil, err
	}
	return results, nil
}
