package graph

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/khanhvtn/netevent-go/graph/generated"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *Resolver) mapUser(m *models.User) (*model.User, error) {
	return &model.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		Roles:     m.Roles,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
func (r *Resolver) mapEvent(m *models.Event) (*model.Event, error) {
	var customizeFields []*model.CustomizeField
	for _, value := range m.CustomizeFields {
		customizeFields = append(customizeFields, &model.CustomizeField{
			Name:     value.Name,
			Type:     value.Type,
			Value:    value.Values,
			Required: value.Required,
		})
	}
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
		CustomizeFields:       customizeFields,
	}, nil
}

func (r *Resolver) mapEventType(m *models.EventType) (*model.EventType, error) {
	return &model.EventType{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Name:      m.Name,
		IsDeleted: m.IsDeleted,
	}, nil
}
func (r *Resolver) mapFacility(m *models.Facility) (*model.Facility, error) {
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
func (r *Resolver) mapFacilityHistory(m *models.FacilityHistory) (*model.FacilityHistory, error) {
	facilityService := r.di.Container.Get(services.FacilityServiceName).(*services.FacilityService)
	eventService := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	facility, err := facilityService.GetOne(bson.M{"_id": m.Facility})
	if err != nil {
		return nil, err
	}
	event, err := eventService.GetOne(bson.M{"_id": m.Event})
	if err != nil {
		return nil, err
	}
	graphModelFacility, err := r.mapFacility(facility)
	if err != nil {
		return nil, err
	}
	graphModelEvent, err := r.mapEvent(event)
	if err != nil {
		return nil, err
	}
	return &model.FacilityHistory{
		ID:         m.ID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		BorrowDate: m.BorrowDate,
		ReturnDate: m.ReturnDate,
		Event:      graphModelEvent,
		Facility:   graphModelFacility,
	}, nil
}
func (r *Resolver) mapParticipant(m *models.Participant) (*model.Participant, error) {
	eventService := r.di.Container.Get(services.EventServiceName).(*services.EventService)

	event, err := eventService.GetOne(bson.M{"_id": m.Event})
	if err != nil {
		return nil, err
	}
	graphModelEvent, err := r.mapEvent(event)
	if err != nil {
		return nil, err
	}
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
		Event:                graphModelEvent,
	}, nil
}
func (r *Resolver) mapTask(m *models.Task) (*model.Task, error) {
	userService := r.di.Container.Get(services.UserServiceName).(*services.UserService)
	eventService := r.di.Container.Get(services.EventServiceName).(*services.EventService)
	user, err := userService.GetOne(bson.M{"_id": m.User})
	if err != nil {
		return nil, err
	}
	event, err := eventService.GetOne(bson.M{"_id": m.Event})
	if err != nil {
		return nil, err
	}
	graphModelUser, err := r.mapUser(user)
	if err != nil {
		return nil, err
	}
	graphModelEvent, err := r.mapEvent(event)
	if err != nil {
		return nil, err
	}
	return &model.Task{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Name:      m.Name,
		Type:      m.Type,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
		Event:     graphModelEvent,
		User:      graphModelUser,
	}, nil
}

func (r *Resolver) GetUserFromContext(ctx context.Context, service *services.UserService) (*models.User, error) {
	//get gin context
	ginContext := ctx.Value("gincontext").(*gin.Context)
	encryptedCookie, err := ginContext.Cookie("netevent")
	if err != nil {
		return nil, errors.New("access denied")
	}
	//decrypt cookie
	id, err := utilities.Decrypted([]byte(encryptedCookie))
	if err != nil {
		return nil, err
	}
	objectId, err := utilities.ConvertStringIdToObjectID(string(id))
	if err != nil {
		return nil, err
	}
	//get user based specific id
	user, err := service.GetOne(bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}
	return user, nil
}
