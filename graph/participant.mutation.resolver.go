package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateParticipant(ctx context.Context, input model.NewParticipant) (*model.Participant, error) {
	service := r.di.Container.Get(services.ParticipantServiceName).(*services.ParticipantService)
	//check input
	if err := service.ValidateNewParticipant(input); err != nil {
		return nil, err
	}
	newParticipant, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapParticipant(newParticipant)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) UpdateParticipant(ctx context.Context, id string, input model.UpdateParticipant) (*model.Participant, error) {
	service := r.di.Container.Get(services.ParticipantServiceName).(*services.ParticipantService)
	//check input
	if err := service.ValidateUpdateParticipant(id, input); err != nil {
		return nil, err
	}
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	updatedParticipant, err := service.UpdateOne(bson.M{"_id": objectId}, input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapParticipant(updatedParticipant)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) DeleteParticipant(ctx context.Context, id string) (*model.Participant, error) {
	service := r.di.Container.Get(services.ParticipantServiceName).(*services.ParticipantService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	deletedParticipant, err := service.DeleteOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	results, err := r.mapParticipant(deletedParticipant)
	if err != nil {
		return nil, err
	}
	return results, nil
}
