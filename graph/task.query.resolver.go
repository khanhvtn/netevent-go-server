package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *taskResolver) Event(ctx context.Context, obj *model.Task) (*model.Event, error) {
	service := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	event, err := service.GetOne(bson.M{"_id": obj.Event.ID})
	if err != nil {
		return nil, err
	}
	results, err := r.mapEvent(event)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *taskResolver) User(ctx context.Context, obj *model.Task) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	user, err := service.GetOne(bson.M{"_id": obj.User.ID})
	if err != nil {
		return nil, err
	}
	results, err := r.mapUser(user)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *queryResolver) Tasks(ctx context.Context) ([]*model.Task, error) {
	return nil, nil
}
func (r *queryResolver) Task(ctx context.Context, id string) (*model.Task, error) {
	return nil, nil
}
