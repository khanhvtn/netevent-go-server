package graph

import (
	"context"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	service := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	//check input
	if err := service.ValidateNewTask(input); err != nil {
		if err := r.extractErrs(err, ctx); err != nil {
			return nil, err
		}
		return nil, nil
	}
	newTask, err := service.Create(input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapTask(newTask)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) UpdateTask(ctx context.Context, id string, input model.UpdateTask) (*model.Task, error) {
	service := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	//check input
	if err := service.ValidateUpdateTask(id, input); err != nil {
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
	updatedTask, err := service.UpdateOne(bson.M{"_id": objectId}, input)
	if err != nil {
		return nil, err
	}
	results, err := r.mapTask(updatedTask)
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (r *mutationResolver) DeleteTask(ctx context.Context, id string) (*model.Task, error) {
	service := r.di.Container.Get(services.TaskServiceName).(*services.TaskService)
	//convert string id to object id
	objectId, err := utilities.ConvertStringIdToObjectID(id)
	if err != nil {
		return nil, err
	}
	deletedTask, err := service.DeleteOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	results, err := r.mapTask(deletedTask)
	if err != nil {
		return nil, err
	}
	return results, nil
}
