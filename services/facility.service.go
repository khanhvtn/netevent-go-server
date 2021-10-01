package services

import (
	"errors"
	"unsafe"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/helpers"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var FacilityServiceName = "FacilityServiceName"

type FacilityService struct {
	FacilityRepository *FacilityRepository
}

/* GetAll: get all data based on condition*/
func (u *FacilityService) GetAll(filter model.FacilityFilter) ([]*models.Facility, error) {
	//setup filter field
	finalFilter := bson.M{}
	filterByString := make([]bson.M, 0)
	opts := options.Find()

	//set filter name
	var keySearch string = ""
	if filter.DefaultFilter.Search != nil {
		keySearch = *filter.DefaultFilter.Search
	}
	filterByString = append(filterByString, bson.M{
		"name": bson.M{"$regex": primitive.Regex{Pattern: keySearch, Options: "i"}},
	})
	filterByString = append(filterByString, bson.M{
		"code": bson.M{"$regex": primitive.Regex{Pattern: keySearch, Options: "i"}},
	})
	filterByString = append(filterByString, bson.M{
		"type": bson.M{"$regex": primitive.Regex{Pattern: keySearch, Options: "i"}},
	})
	//set date filter
	//for status
	var status bool = false
	if filter.Status != nil {
		status = *filter.Status
	}
	finalFilter["status"] = status

	//for updatedAt
	if filter.DefaultFilter.UpdatedAtDateFrom != nil && filter.DefaultFilter.UpdatedAtDateTo != nil {
		finalFilter["updatedAt"] = bson.M{
			"$gte": filter.DefaultFilter.UpdatedAtDateFrom,
			"$lte": filter.DefaultFilter.UpdatedAtDateTo,
		}
	}
	//for createdAt
	if filter.DefaultFilter.CreatedAtDateFrom != nil && filter.DefaultFilter.CreatedAtDateTo != nil {
		finalFilter["createdAt"] = bson.M{
			"$gte": filter.DefaultFilter.CreatedAtDateFrom,
			"$lte": filter.DefaultFilter.CreatedAtDateTo,
		}
	}

	//set the number of record that will display
	var take int64 = 10 //take 10 records in default
	if filter.DefaultFilter.Take != nil {
		take = *(*int64)(unsafe.Pointer(filter.DefaultFilter.Take))
	}
	opts.SetLimit(take)
	//set paging
	var page int64 = 1 //target page 1 in default
	if filter.DefaultFilter.Page != nil {
		page = *(*int64)(unsafe.Pointer(filter.DefaultFilter.Page))
	}
	opts.SetSkip((page - 1) * take)

	//set isDeleted filter
	var isDeleted = false
	if filter.DefaultFilter.IsDeleted != nil {
		isDeleted = *filter.DefaultFilter.IsDeleted
	}
	finalFilter["isDeleted"] = isDeleted

	//set filter for string field
	finalFilter["$or"] = filterByString

	return u.FacilityRepository.FindAll(finalFilter, opts)
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

//validation
func (u *FacilityService) ValidateNewFacility(newFacility model.NewFacility) error {
	return validation.ValidateStruct(&newFacility,
		validation.Field(&newFacility.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			event, err := u.GetOne(bson.M{"name": name.(string)})
			if _, ok := err.(*helpers.ErrNotFound); err != nil && !ok {
				return err
			}
			if event != nil {
				return errors.New("name already existed")
			}
			return nil

		})),
		validation.Field(&newFacility.Code, validation.Required.Error("code must not be blanked")),
		validation.Field(&newFacility.Type, validation.Required.Error("type must not be blanked")),
	)
}

func (u *FacilityService) ValidateUpdateFacility(id string, updateFacility model.UpdateFacility) error {
	return validation.ValidateStruct(&updateFacility,
		validation.Field(&updateFacility.Name, validation.Required.Error("name must not be blanked"), validation.By(func(name interface{}) error {
			//convert id string to object id
			objectId, err := utilities.ConvertStringIdToObjectID(id)
			if err != nil {
				return err
			}
			//get current facility
			currentFacility, err := u.GetOne(bson.M{"_id": objectId})
			if err != nil {
				return err
			}
			//check email existed or not
			facility, err := u.GetOne(bson.M{"name": name.(string)})
			if err != nil {
				return err
			}
			if facility != nil && facility.Name != currentFacility.Name {
				return errors.New("name already existed")
			}
			return nil

		})),
		validation.Field(&updateFacility.Code, validation.Required.Error("code must not be blanked")),
		validation.Field(&updateFacility.Type, validation.Required.Error("type must not be blanked")),
		validation.Field(&updateFacility.IsDeleted, validation.Required.Error("delete status must not be blanked")),
		validation.Field(&updateFacility.Status, validation.Required.Error("status must not be blanked")),
	)
}
