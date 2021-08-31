package services

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var FacilityHistoryServiceName = "FacilityHistoryServiceName"

type FacilityHistoryService struct {
	FacilityHistoryRepository *FacilityHistoryRepository
}

/* GetAll: get all data based on condition*/
func (u *FacilityHistoryService) GetAll(condition bson.M) ([]*models.FacilityHistory, error) {
	return u.FacilityHistoryRepository.FindAll(condition)
}

/*GetOne: get one record from a collection  */
func (u *FacilityHistoryService) GetOne(filter bson.M) (*models.FacilityHistory, error) {
	return u.FacilityHistoryRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *FacilityHistoryService) Create(newFacilityHistory model.NewFacilityHistory) (*models.FacilityHistory, error) {
	//get facility, event
	facilityId, err := primitive.ObjectIDFromHex(newFacilityHistory.FacilityID)
	if err != nil {
		return nil, err
	}
	eventId, err := primitive.ObjectIDFromHex(*newFacilityHistory.EventID)
	if err != nil {
		return nil, err
	}

	//convert to bson.M
	currentTime := time.Now()
	facilityHistory := models.FacilityHistory{
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
		Facility:   facilityId,
		BorrowDate: newFacilityHistory.BorrowDate,
		ReturnDate: newFacilityHistory.ReturnDate,
		Event:      eventId,
	}
	return u.FacilityHistoryRepository.Create(&facilityHistory)
}

/*UpdateOne: update one record from a collection*/
func (u FacilityHistoryService) UpdateOne(filter bson.M, update model.UpdateFacilityHistory) (*models.FacilityHistory, error) {
	//convert interface to bson
	bsonUpdate, err := utilities.InterfaceToBsonM(update)
	if err != nil {
		return nil, err
	}
	return u.FacilityHistoryRepository.UpdateOne(filter, bsonUpdate)
}

//DeleteOne func is to update one record from a collection
func (u FacilityHistoryService) DeleteOne(filter bson.M) (*models.FacilityHistory, error) {
	return u.FacilityHistoryRepository.DeleteOne(filter)
}

//validation
func (u *FacilityHistoryService) ValidateNewFacilityHistory(newFacilityHistory model.NewFacilityHistory) error {
	return validation.ValidateStruct(&newFacilityHistory,
		validation.Field(&newFacilityHistory.FacilityID, validation.Required.Error("facility id must not be blanked"), validation.By(func(id interface{}) error {
			if ok := primitive.IsValidObjectID(id.(string)); !ok {
				return errors.New("invalid id")
			}
			return nil

		})),
		validation.Field(&newFacilityHistory.BorrowDate, validation.Required.Error("Borrow date must not be blanked")),
		validation.Field(&newFacilityHistory.ReturnDate, validation.Required.Error("Return date password must not be blanked")),
	)
}

func (u *FacilityHistoryService) ValidateUpdateFacilityHistory(id string, updateFacilityHistory model.UpdateFacilityHistory) error {
	return validation.ValidateStruct(&updateFacilityHistory,
		validation.Field(&updateFacilityHistory.FacilityID, validation.Required.Error("facility id must not be blanked"), validation.By(func(id interface{}) error {
			if ok := primitive.IsValidObjectID(id.(string)); !ok {
				return errors.New("invalid id")
			}
			return nil

		})),
		validation.Field(&updateFacilityHistory.BorrowDate, validation.Required.Error("Borrow date must not be blanked")),
		validation.Field(&updateFacilityHistory.ReturnDate, validation.Required.Error("Return date password must not be blanked")),
	)
}
