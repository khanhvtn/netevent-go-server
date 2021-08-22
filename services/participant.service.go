package services

import (
	"time"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ParticipantServiceName = "ParticipantServiceName"

type ParticipantService struct {
	ParticipantRepository *ParticipantRepository
}

/* GetAll: get all data based on condition*/
func (u *ParticipantService) GetAll(condition *bson.M) ([]*models.Participant, error) {
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
func (u ParticipantService) UpdateOne(filter bson.M, update bson.M) (*models.Participant, error) {
	return u.ParticipantRepository.UpdateOne(filter, update)
}

//DeleteOne func is to update one record from a collection
func (u ParticipantService) DeleteOne(filter bson.M) (*models.Participant, error) {
	return u.ParticipantRepository.DeleteOne(filter)
}
