package services

import (
	"time"

	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var FacilityHistoryServiceName = "FacilityHistoryServiceName"

type FacilityHistoryService struct {
	FacilityHistoryRepository *FacilityHistoryRepository
}

/* GetAll: get all data based on condition*/
func (u *FacilityHistoryService) GetAll(condition *bson.M) ([]*models.FacilityHistory, error) {
	return u.FacilityHistoryRepository.FindAll(condition)
}

/*GetOne: get one record from a collection  */
func (u *FacilityHistoryService) GetOne(filter bson.M) (*models.FacilityHistory, error) {
	return u.FacilityHistoryRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *FacilityHistoryService) Create(newFacilityHistory *model.NewFacilityHistory) (*models.FacilityHistory, error) {

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
func (u FacilityHistoryService) UpdateOne(filter bson.M, update bson.M) (*models.FacilityHistory, error) {
	return u.FacilityHistoryRepository.UpdateOne(filter, update)
}

//DeleteOne func is to update one record from a collection
func (u FacilityHistoryService) DeleteOne(filter bson.M) (*models.FacilityHistory, error) {
	return u.FacilityHistoryRepository.DeleteOne(filter)
}
