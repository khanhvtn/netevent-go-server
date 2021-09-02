package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *participantResolver) Event(ctx context.Context, obj *model.Participant) (*model.Event, error) {
	eventService := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	participantService := r.di.Container.Get(services.ParticipantServiceName).(*services.ParticipantService)
	participant, err := participantService.GetOne(bson.M{"_id": obj.ID})
	if err != nil {
		return nil, err
	}
	event, err := eventService.GetOne(bson.M{"_id": participant.Event})
	if err != nil {
		return nil, err
	}
	results, err := r.mapEvent(event)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *queryResolver) Participants(ctx context.Context) ([]*model.Participant, error) {
	service := r.di.Container.Get(services.ParticipantServiceName).(*services.ParticipantService)
	participants, err := service.GetAll(bson.M{})
	if err != nil {
		return nil, err
	}
	results := make([]*model.Participant, 0)
	for _, participant := range participants {
		mappedParticipant, err := r.mapParticipant(participant)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedParticipant)
	}
	return results, nil
}

func (r *queryResolver) Participant(ctx context.Context, id string) (*model.Participant, error) {
	service := r.di.Container.Get(services.ParticipantServiceName).(*services.ParticipantService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	//get participant based specific id
	participant, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapParticipant(participant)
	if err != nil {
		return nil, err
	}
	return result, nil
}
