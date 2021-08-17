package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *participantResolver) Event(ctx context.Context, obj *model.Participant) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	event, err := service.GetOne(bson.M{"_id": obj.Event.ID.Hex()})
	if err != nil {
		return nil, err
	}
	results, _ := mapEvent(event)
	return results, nil
}
