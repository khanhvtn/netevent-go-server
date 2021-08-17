package graph

import (
	"github.com/khanhvtn/netevent-go/graph/generated"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	di *services.DI
}

func Init(di *services.DI) *Resolver {
	return &Resolver{
		di: di,
	}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type eventResolver struct{ *Resolver }
type facilityHistoryResolver struct{ *Resolver }
type participantResolver struct{ *Resolver }
type taskResolver struct{ *Resolver }

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// FacilityHistory returns generated.FacilityHistoryResolver implementation.
func (r *Resolver) FacilityHistory() generated.FacilityHistoryResolver {
	return &facilityHistoryResolver{r}
}

// Participant returns generated.ParticipantResolver implementation.
func (r *Resolver) Participant() generated.ParticipantResolver { return &participantResolver{r} }

// Task returns generated.TaskResolver implementation.
func (r *Resolver) Task() generated.TaskResolver { return &taskResolver{r} }

func mapUser(m *models.User) (*model.User, error) {
	return &model.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		Roles:     m.Roles,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
func mapEvent(m *models.Event) (*model.Event, error) {
	return &model.Event{
		ID:                    m.ID,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		Tags:                  m.Tags,
		IsApproved:            m.IsApproved,
		IsFinished:            m.IsFinished,
		Name:                  m.Name,
		Language:              m.Language,
		Mode:                  m.Mode,
		Location:              m.Location,
		Accommodation:         m.Accommodation,
		RegistrationCloseDate: m.RegistrationCloseDate,
		StartDate:             m.StartDate,
		EndDate:               m.EndDate,
		MaxParticipants:       m.MaxParticipants,
		Description:           m.Description,
		Budget:                m.Budget,
		Image:                 m.Image,
		IsDeleted:             m.IsDeleted,
	}, nil
}

func mapEventType(m *models.EventType) (*model.EventType, error) {
	return &model.EventType{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Name:      m.Name,
		IsDeleted: m.IsDeleted,
	}, nil
}
func mapFacility(m *models.Facility) (*model.Facility, error) {
	return &model.Facility{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Name:      m.Name,
		Status:    m.Status,
		Code:      m.Code,
		Type:      m.Type,
		IsDeleted: m.IsDeleted,
	}, nil
}
func mapFacilityHistory(m *models.FacilityHistory) (*model.FacilityHistory, error) {
	return &model.FacilityHistory{
		ID:         m.ID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		BorrowDate: m.BorrowDate,
		ReturnDate: m.ReturnDate,
	}, nil
}
func mapParticipant(m *models.Participant) (*model.Participant, error) {
	return &model.Participant{
		ID:                   m.ID,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
		IsValid:              m.IsValid,
		IsAttended:           m.IsAttended,
		Email:                m.Email,
		Name:                 m.Name,
		Academic:             m.Academic,
		School:               m.School,
		Major:                m.Major,
		Phone:                m.Phone,
		Dob:                  m.DOB,
		ExpectedGraduateDate: m.ExpectedGraduateDate,
	}, nil
}
func mapTask(m *models.Task) (*model.Task, error) {
	return &model.Task{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Name:      m.Name,
		Type:      m.Type,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
	}, nil
}
