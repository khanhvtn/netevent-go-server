package services

import (
	"errors"
	"math"
	"unsafe"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/khanhvtn/netevent-go/graph/model"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserServiceName = "UserServiceName"

// UserService handles the creation, modification and deletion of users.
type UserService struct {
	UserRepository *UserRepository
}

/* GetAll: get all data based on condition*/
func (u *UserService) GetAll(filter model.UserFilter) ([]*models.User, *model.PageInfo, error) {
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
		"email": bson.M{"$regex": primitive.Regex{Pattern: keySearch, Options: "i"}},
	})

	//set roles filter
	if filter.Roles != nil {
		finalFilter["roles"] = bson.M{
			"$all": filter.Roles,
		}
	}

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

	//set filter for string field
	finalFilter["$or"] = filterByString

	//get total number page
	errChan := make(chan error, 2)
	totalPageChan := make(chan int)
	var totalPage *int = nil

	go func(totalPageChan chan<- int, errChan chan<- error, finalFilter primitive.M, opts *options.FindOptions, take int64) {
		if users, err := u.UserRepository.FindAll(finalFilter, opts); err != nil {
			errChan <- err
			close(totalPageChan)
		} else {
			totalUsers := len(users)
			totalPage := int(math.Ceil(float64(totalUsers) / float64(take)))
			errChan <- nil
			totalPageChan <- totalPage
		}

	}(totalPageChan, errChan, finalFilter, opts, take)

	//set paging
	var page int64 = 1 //target page 1 in default
	if filter.DefaultFilter.Page != nil {
		page = *(*int64)(unsafe.Pointer(filter.DefaultFilter.Page))
	}
	opts.SetSkip((page - 1) * take)

	//get users
	var users []*models.User = nil
	usersChan := make(chan []*models.User)
	go func(usersChan chan<- []*models.User, errChan chan<- error, finalFilter primitive.M, opts *options.FindOptions) {
		if users, err := u.UserRepository.FindAll(finalFilter, opts); err != nil {
			errChan <- err
			close(usersChan)
		} else {
			errChan <- nil
			usersChan <- users
		}

	}(usersChan, errChan, finalFilter, opts)

	if totalPageValue, ok := <-totalPageChan; ok {
		totalPage = &totalPageValue
	}
	if usersValue, ok := <-usersChan; ok {
		users = usersValue
	}

	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, nil, err
		}
	}

	return users, &model.PageInfo{TotalPage: *totalPage, CurrentPage: int(page)}, nil
}

/*GetOne: get one record from a collection  */
func (u *UserService) GetOne(filter bson.M) (*models.User, error) {
	return u.UserRepository.FindOne(filter)
}

/*Create: create a new record to a collection*/
func (u *UserService) Create(newUser model.NewUser) (*models.User, error) {
	return u.UserRepository.Create(newUser)
}

/*UpdateOne: update one record from a collection*/
func (u UserService) UpdateOne(filter bson.M, update bson.M) (*models.User, error) {
	return u.UserRepository.UpdateOne(filter, update)
}

//DeleteOne func is to update one record from a collection
func (u UserService) DeleteOne(filter bson.M) (*models.User, error) {
	return u.UserRepository.DeleteOne(filter)
}

func (u UserService) Login(input model.Login) (*models.User, error) {
	user, err := u.GetOne(bson.M{"email": input.Email})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid user or password")
	}
	//check password
	if ok := utilities.CheckPasswordHash(input.Password, user.Password); !ok {
		return nil, errors.New("invalid user or password")
	}
	return user, nil
}

//validation
func (u *UserService) ValidateNewUser(newUser model.NewUser) error {
	return validation.ValidateStruct(&newUser,
		validation.Field(&newUser.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email"), validation.By(func(email interface{}) error {
			user, err := u.GetOne(bson.M{"email": email.(string)})
			if err != nil {
				return err
			}
			if user != nil {
				return errors.New("email already existed")
			}
			return nil

		})),
		validation.Field(&newUser.Password, validation.Required.Error("password must not be blanked")),
		validation.Field(&newUser.ConfirmPassword, validation.Required.Error("confirm password must not be blanked"), validation.In(newUser.Password).Error("confirm password must be identical with Password")),
		validation.Field(&newUser.Roles, validation.Required.Error("Role must not be blanked")),
	)
}

func (u *UserService) ValidateUpdateUser(id string, updateUser model.UpdateUser) error {
	return validation.ValidateStruct(&updateUser,
		validation.Field(&updateUser.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email"), validation.By(func(email interface{}) error {
			//get current user
			currentUser, err := u.GetOne(bson.M{"_id": id})
			if err != nil {
				return err
			}
			//check email existed or not
			user, err := u.GetOne(bson.M{"email": email.(string)})
			if err != nil {
				return err
			}
			if user != nil && user.Email != currentUser.Email {
				return errors.New("email already existed")
			}
			return nil

		})),
		validation.Field(&updateUser.Password, validation.Required.Error("password must not be blanked")),
		validation.Field(&updateUser.Roles, validation.Required.Error("Role must not be blanked")),
	)
}
func (u *UserService) ValidateLogin(login model.Login) error {
	return validation.ValidateStruct(&login,
		validation.Field(&login.Email, validation.Required.Error("email must not be blanked"), is.Email.Error("invalid email")),
		validation.Field(&login.Password, validation.Required.Error("password must not be blanked")),
	)
}

func (u *UserService) HashPassword(newUser *model.NewUser) error {
	hashPassword, err := utilities.HashPassword(newUser.Password)
	if err != nil {
		return err
	}
	newUser.Password = hashPassword
	return nil
}
func (u *UserService) HashPasswordUpdateUser(updateUser *model.UpdateUser) error {
	hashPassword, err := utilities.HashPassword(updateUser.Password)
	if err != nil {
		return err
	}
	updateUser.Password = hashPassword
	return nil
}
