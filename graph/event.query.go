package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *eventResolver) Tasks(ctx context.Context, obj *model.Event) ([]*model.Task, error) {
	service := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	tasks, err := service.GetAll(&bson.M{"event": obj.ID.Hex()})
	if err != nil {
		return nil, err
	}
	results := make([]*model.Task, 0)
	for _, task := range tasks {
		mappedTask, _ := mapTask(task)
		results = append(results, mappedTask)
	}
	return results, nil
}

func (r *eventResolver) FacilityHistories(ctx context.Context, obj *model.Event) ([]*model.FacilityHistory, error) {
	service := r.di.Container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	facilityHistories, err := service.GetAll(&bson.M{"event": obj.ID.Hex()})
	if err != nil {
		return nil, err
	}
	results := make([]*model.FacilityHistory, 0)
	for _, facilityHistory := range facilityHistories {
		mappedFacilityHistory, _ := mapFacilityHistory(facilityHistory)
		results = append(results, mappedFacilityHistory)
	}
	return results, nil
}

func (r *eventResolver) EventType(ctx context.Context, obj *model.Event) (*model.EventType, error) {
	service := r.di.Container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	eventType, err := service.GetOne(bson.M{"_id": obj.EventType.ID.Hex()})
	if err != nil {
		return nil, err
	}
	results, _ := mapEventType(eventType)
	return results, nil
}

func (r *eventResolver) Owner(ctx context.Context, obj *model.Event) (*model.User, error) {
	service := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	owner, err := service.GetOne(bson.M{"_id": obj.Owner.ID.Hex()})
	if err != nil {
		return nil, err
	}
	results, _ := mapUser(owner)
	return results, nil
}
