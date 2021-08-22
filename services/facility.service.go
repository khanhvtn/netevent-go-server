package services

import (
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"go.mongodb.org/mongo-driver/bson"
)

var FacilityServiceName = "FacilityServiceName"

type FacilityService struct {
	FacilityRepository *FacilityRepository
}

/* GetAll: get all data based on condition*/
func (u *FacilityService) GetAll(condition *bson.M) ([]*models.Facility, error) {
	return u.FacilityRepository.FindAll(condition)
}

/*GetOne: get one record from a collection  */
func (u *FacilityService) GetOne(filter bson.M) (*models.Facility, error) {
	return u.FacilityRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *FacilityService) Create(newFacility model.NewFacility) (*models.Facility, error) {
	return u.FacilityRepository.Create(newFacility)
}

/*UpdateOne: update one record from a collection*/
func (u FacilityService) UpdateOne(filter bson.M, update bson.M) (*models.Facility, error) {
	return u.FacilityRepository.UpdateOne(filter, update)
}

//DeleteOne func is to update one record from a collection
func (u FacilityService) DeleteOne(filter bson.M) (*models.Facility, error) {
	return u.FacilityRepository.DeleteOne(filter)
}
