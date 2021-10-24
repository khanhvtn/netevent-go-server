package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateFacilityHistory(ctx context.Context, input model.NewFacilityHistory) (*model.FacilityHistory, error) {
	service := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	//check input
	if err := service.ValidateNewFacilityHistory(input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	newFacilityHistory, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapFacilityHistory(newFacilityHistory)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) UpdateFacilityHistory(ctx context.Context, id string, input model.UpdateFacilityHistory) (*model.FacilityHistory, error) {
	service := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	//check input
	if err := service.ValidateUpdateFacilityHistory(id, input); err != nil {
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
	updatedFacilityHistory, err := service.UpdateOne(bson.M{"_id": objectId}, input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapFacilityHistory(updatedFacilityHistory)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) DeleteFacilityHistory(ctx context.Context, id string) (*model.FacilityHistory, error) {
	service := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	deletedFacilityHistory, err := service.DeleteOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	results, err := r.mapFacilityHistory(deletedFacilityHistory)
	if err != nil {
		return nil, err
	}
	return results, nil
}
