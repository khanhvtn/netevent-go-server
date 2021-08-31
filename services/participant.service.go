package services

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ParticipantServiceName = "ParticipantServiceName"

type ParticipantService struct {
	ParticipantRepository *ParticipantRepository
}

/* GetAll: get all data based on condition*/
func (u *ParticipantService) GetAll(condition bson.M) ([]*models.Participant, error) {
	return u.ParticipantRepository.FindAll(condition)
}

/*GetOne: get one record from a collection  */
func (u *ParticipantService) GetOne(filter bson.M) (*models.Participant, error) {
	return u.ParticipantRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *ParticipantService) Create(newParticipant model.NewParticipant) (*models.Participant, error) {

	//get event
	eventId, err := primitive.ObjectIDFromHex(newParticipant.EventID)
	if err != nil {
		return nil, err
	}

	//convert to bson.M
	currentTime := time.Now()
	participant := models.Participant{
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		IsValid:              false,
		IsAttended:           false,
		Event:                eventId,
		Email:                newParticipant.Email,
		Name:                 newParticipant.Name,
		Academic:             newParticipant.Academic,
		School:               newParticipant.School,
		Major:                newParticipant.Major,
		Phone:                newParticipant.Phone,
		DOB:                  newParticipant.Dob,
		ExpectedGraduateDate: newParticipant.ExpectedGraduateDate,
	}
	return u.ParticipantRepository.Create(participant)
}

/*UpdateOne: update one record from a collection*/
func (u ParticipantService) UpdateOne(filter bson.M, update model.UpdateParticipant) (*models.Participant, error) {
	bsonUpdate, err := utilities.InterfaceToBsonM(update)
	if err != nil {
		return nil, err
	}
	return u.ParticipantRepository.UpdateOne(filter, bsonUpdate)
}

//DeleteOne func is to update one record from a collection
func (u ParticipantService) DeleteOne(filter bson.M) (*models.Participant, error) {
	return u.ParticipantRepository.DeleteOne(filter)
}

//validation
func (u *ParticipantService) ValidateNewParticipant(newParticipant model.NewParticipant) error {
	return validation.ValidateStruct(&newParticipant,
		validation.Field(&newParticipant.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email"), validation.By(func(email interface{}) error {
			participant, err := u.GetOne(bson.M{"email": email.(string)})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			if participant != nil {
				return errors.New("email already existed")
			}
			return nil

		})),
		validation.Field(&newParticipant.EventID, validation.Required.Error("event id must not be blanked")),
		validation.Field(&newParticipant.Name, validation.Required.Error("name must not be blanked")),
		validation.Field(&newParticipant.Academic, validation.Required.Error("academic must not be blanked")),
		validation.Field(&newParticipant.Dob, validation.Required.Error("date of birth must not be blanked")),
		validation.Field(&newParticipant.ExpectedGraduateDate, validation.Required.Error("expected graduate date must not be blanked")),
		validation.Field(&newParticipant.Major, validation.Required.Error("major must not be blanked")),
		validation.Field(&newParticipant.Phone, validation.Required.Error("phone must not be blanked")),
		validation.Field(&newParticipant.School, validation.Required.Error("school must not be blanked")),
	)
}

func (u *ParticipantService) ValidateUpdateParticipant(id string, updateParticipant model.UpdateParticipant) error {
	return validation.ValidateStruct(&updateParticipant,
		validation.Field(&updateParticipant.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email"), validation.By(func(email interface{}) error {
			//convert string id to object id
			objectId, err := utilities.ConvertStringIdToObjectID(id)
			if err != nil {
				return err
			}
			//get current participant
			currentParticipant, err := u.GetOne(bson.M{"_id": objectId})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			//check email existed or not
			if participant, err := u.GetOne(bson.M{"email": email.(string)}); err != nil {
				if _, ok := err.(*helpers.ErrNotFound); ok {
					return nil
				} else {
					return err
				}
			} else {
				if participant.Name != currentParticipant.Name {
					return errors.New("email already existed")
				} else {
					return nil
				}
			}

		})),
		validation.Field(&updateParticipant.EventID, validation.Required.Error("event id must not be blanked")),
		validation.Field(&updateParticipant.Name, validation.Required.Error("name must not be blanked")),
		validation.Field(&updateParticipant.Academic, validation.Required.Error("academic must not be blanked")),
		validation.Field(&updateParticipant.Dob, validation.Required.Error("date of birth must not be blanked")),
		validation.Field(&updateParticipant.ExpectedGraduateDate, validation.Required.Error("expected graduate date must not be blanked")),
		validation.Field(&updateParticipant.Major, validation.Required.Error("major must not be blanked")),
		validation.Field(&updateParticipant.Phone, validation.Required.Error("phone must not be blanked")),
		validation.Field(&updateParticipant.School, validation.Required.Error("school must not be blanked")),
	)
}
