package graph

import (
	"github.com/khanhvtn/netevent-go/models"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userModel            *models.User
	eventModel           *models.Event
	eventTypeModel       *models.EventType
	facilityModel        *models.Facility
	facilityHistoryModel *models.FacilityHistory
	participantModel     *models.Participant
	taskModel            *models.Task
}

func Init() *Resolver {
	return &Resolver{
		userModel:            &models.User{},
		eventModel:           &models.Event{},
		eventTypeModel:       &models.EventType{},
		facilityModel:        &models.Facility{},
		facilityHistoryModel: &models.FacilityHistory{},
		participantModel:     &models.Participant{},
		taskModel:            &models.Task{},
	}
}
