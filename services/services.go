package services

import (
	"context"
	"time"

	"github.com/khanhvtn/netevent-go/database"
	"github.com/sarulabs/di"
)

type DI struct {
	Container di.Container
}

func New() (*DI, error) {
	// Create the app container.
	// Do not forget to delete it at the end.
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	err = builder.Add(services...)
	if err != nil {
		return nil, err
	}

	app := builder.Build()
	return &DI{Container: app}, nil
}

// Services contains the definitions of the application services.
var services = []di.Def{
	{
		Name: database.MongoCNName,
		Build: func(ctn di.Container) (interface{}, error) {
			return database.ConnectDB(), nil
		},
		Close: func(obj interface{}) error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			return obj.(*database.MongoInstance).Client.Disconnect(ctx)
		},
	},
	{
		Name: UserRepositoryName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &UserRepository{
				MongoCN: ctn.Get(database.MongoCNName).(*database.MongoInstance),
			}, nil
		},
	},
	{
		Name: UserServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &UserService{
				UserRepository: ctn.Get(UserRepositoryName).(*UserRepository),
			}, nil
		},
	},
	{
		Name: EventTypeRepositoryName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &EventTypeRepository{
				MongoCN: ctn.Get(database.MongoCNName).(*database.MongoInstance),
			}, nil
		},
	},
	{
		Name: EventTypeServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &EventTypeService{
				EventTypeRepository: ctn.Get(EventTypeRepositoryName).(*EventTypeRepository),
			}, nil
		},
	},
	{
		Name: FacilityRepositoryName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &FacilityRepository{
				MongoCN: ctn.Get(database.MongoCNName).(*database.MongoInstance),
			}, nil
		},
	},
	{
		Name: FacilityServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &FacilityService{
				FacilityRepository: ctn.Get(FacilityRepositoryName).(*FacilityRepository),
			}, nil
		},
	},
	{
		Name: EventRepositoryName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &EventRepository{
				MongoCN: ctn.Get(database.MongoCNName).(*database.MongoInstance),
			}, nil
		},
	},
	{
		Name: EventServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &EventService{
				EventRepository:           ctn.Get(EventRepositoryName).(*EventRepository),
				TaskRepository:            ctn.Get(TaskRepositoryName).(*TaskRepository),
				FacilityHistoryRepository: ctn.Get(FacilityHistoryRepositoryName).(*FacilityHistoryRepository),
			}, nil
		},
	},
	{
		Name: FacilityHistoryRepositoryName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &FacilityHistoryRepository{
				MongoCN: ctn.Get(database.MongoCNName).(*database.MongoInstance),
			}, nil
		},
	},
	{
		Name: FacilityHistoryServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &FacilityHistoryService{
				FacilityHistoryRepository: ctn.Get(FacilityHistoryRepositoryName).(*FacilityHistoryRepository),
			}, nil
		},
	},
	{
		Name: ParticipantRepositoryName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &ParticipantRepository{
				MongoCN: ctn.Get(database.MongoCNName).(*database.MongoInstance),
			}, nil
		},
	},
	{
		Name: ParticipantServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &ParticipantService{
				ParticipantRepository: ctn.Get(ParticipantRepositoryName).(*ParticipantRepository),
			}, nil
		},
	},
	{
		Name: TaskRepositoryName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &TaskRepository{
				MongoCN: ctn.Get(database.MongoCNName).(*database.MongoInstance),
			}, nil
		},
	},
	{
		Name: TaskServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &TaskService{
				TaskRepository: ctn.Get(TaskRepositoryName).(*TaskRepository),
			}, nil
		},
	},
}
