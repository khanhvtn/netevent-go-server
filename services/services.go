package services

import (
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
		Name: UserServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &UserService{
				MongoCN: database.MongoCN,
			}, nil
		},
	},
	{
		Name: EventServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &EventService{
				MongoCN:                database.MongoCN,
				TaskService:            ctn.Get(TaskServiceName).(*TaskService),
				FacilityHistoryService: ctn.Get(FacilityHistoryServiceName).(*FacilityHistoryService),
			}, nil
		},
	},
	{
		Name: EventTypeServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &EventTypeService{
				MongoCN: database.MongoCN,
			}, nil
		},
	},
	{
		Name: FacilityServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &FacilityService{
				MongoCN: database.MongoCN,
			}, nil
		},
	},
	{
		Name: FacilityHistoryServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &FacilityHistoryService{
				MongoCN:         database.MongoCN,
				FacilityService: ctn.Get(FacilityServiceName).(*FacilityService),
				EventService:    ctn.Get(EventServiceName).(*EventService),
			}, nil
		},
	},
	{
		Name: ParticipantServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &ParticipantService{
				MongoCN:      database.MongoCN,
				EventService: ctn.Get(EventServiceName).(*EventService),
			}, nil
		},
	},
	{
		Name: TaskServiceName,
		Build: func(ctn di.Container) (interface{}, error) {
			return &TaskService{
				MongoCN:      database.MongoCN,
				EventService: ctn.Get(EventServiceName).(*EventService),
				UserService:  ctn.Get(UserServiceName).(*UserService),
			}, nil
		},
	},
}
