package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *taskResolver) Event(ctx context.Context, obj *model.Task) (*model.Event, error) {
	eventService := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	taskService := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	task, err := taskService.GetOne(bson.M{"_id": obj.ID})
	if err != nil {
		return nil, err
	}
	event, err := eventService.GetOne(bson.M{"_id": task.Event})
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
	userService := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	taskService := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	task, err := taskService.GetOne(bson.M{"_id": obj.ID})
	if err != nil {
		return nil, err
	}
	user, err := userService.GetOne(bson.M{"_id": task.User})
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
	service := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	tasks, err := service.GetAll(bson.M{})
	if err != nil {
		return nil, err
	}
	results := make([]*model.Task, 0)
	for _, participant := range tasks {
		mappedTask, err := r.mapTask(participant)
		if err != nil {
			return nil, err
		}
		results = append(results, mappedTask)
	}
	return results, nil
}
func (r *queryResolver) Task(ctx context.Context, id string) (*model.Task, error) {
	service := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	//get task based specific id
	task, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	result, err := r.mapTask(task)
	if err != nil {
		return nil, err
	}
	return result, nil
}
