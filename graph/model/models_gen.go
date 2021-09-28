// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomizeField struct {
	Name     string   `json:"name" bson:"name"`
	Type     string   `json:"type" bson:"type"`
	Value    []string `json:"value" bson:"value"`
	Required bool     `json:"required" bson:"required"`
}

type Event struct {
	ID                    primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt             time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt             time.Time          `json:"updatedAt" bson:"updatedAt"`
	Tags                  []string           `json:"tags" bson:"tags"`
	IsApproved            bool               `json:"isApproved" bson:"isApproved"`
	Reviewer              *User              `json:"reviewer" bson:"reviewer"`
	IsFinished            bool               `json:"isFinished" bson:"isFinished"`
	Tasks                 []*Task            `json:"tasks" bson:"tasks"`
	FacilityHistories     []*FacilityHistory `json:"facilityHistories" bson:"facilityHistories"`
	Name                  string             `json:"name" bson:"name"`
	Language              string             `json:"language" bson:"language"`
	EventType             *EventType         `json:"eventType" bson:"eventType"`
	Mode                  string             `json:"mode" bson:"mode"`
	Location              string             `json:"location" bson:"location"`
	Accommodation         string             `json:"accommodation" bson:"accommodation"`
	RegistrationCloseDate time.Time          `json:"registrationCloseDate" bson:"registrationCloseDate"`
	StartDate             time.Time          `json:"startDate" bson:"startDate"`
	EndDate               time.Time          `json:"endDate" bson:"endDate"`
	MaxParticipants       int                `json:"maxParticipants" bson:"maxParticipants"`
	Description           string             `json:"description" bson:"description"`
	Owner                 *User              `json:"owner" bson:"owner"`
	Budget                float64            `json:"budget" bson:"budget"`
	Image                 string             `json:"image" bson:"image"`
	IsDeleted             bool               `json:"isDeleted" bson:"isDeleted"`
	CustomizeFields       []*CustomizeField  `json:"customizeFields" bson:"customizeFields"`
}

type EventType struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	Name      string             `json:"name" bson:"name"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
}

type Facility struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	Status    bool               `json:"status" bson:"status"`
	Name      string             `json:"name" bson:"name"`
	Code      string             `json:"code" bson:"code"`
	Type      string             `json:"type" bson:"type"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
}

type FacilityHistory struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
	Facility   *Facility          `json:"facility" bson:"facility"`
	BorrowDate time.Time          `json:"borrowDate" bson:"borrowDate"`
	ReturnDate time.Time          `json:"returnDate" bson:"returnDate"`
	Event      *Event             `json:"event" bson:"event"`
}

type InputCustomizeField struct {
	Name     string   `json:"name" bson:"name"`
	Type     string   `json:"type" bson:"type"`
	Value    []string `json:"value" bson:"value"`
	Required bool     `json:"required" bson:"required"`
}

type Login struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type NewEvent struct {
	Tags                  []string               `json:"tags" bson:"tags"`
	Tasks                 []*NewTask             `json:"tasks" bson:"tasks"`
	FacilityHistories     []*NewFacilityHistory  `json:"facilityHistories" bson:"facilityHistories"`
	Name                  string                 `json:"name" bson:"name"`
	Language              string                 `json:"language" bson:"language"`
	EventTypeID           string                 `json:"eventTypeId" bson:"eventTypeId"`
	Mode                  string                 `json:"mode" bson:"mode"`
	Location              string                 `json:"location" bson:"location"`
	Accommodation         string                 `json:"accommodation" bson:"accommodation"`
	RegistrationCloseDate time.Time              `json:"registrationCloseDate" bson:"registrationCloseDate"`
	StartDate             time.Time              `json:"startDate" bson:"startDate"`
	EndDate               time.Time              `json:"endDate" bson:"endDate"`
	MaxParticipants       int                    `json:"maxParticipants" bson:"maxParticipants"`
	Description           string                 `json:"description" bson:"description"`
	OwnerID               string                 `json:"ownerId" bson:"ownerId"`
	Budget                float64                `json:"budget" bson:"budget"`
	Image                 string                 `json:"image" bson:"image"`
	CustomizeFields       []*InputCustomizeField `json:"customizeFields" bson:"customizeFields"`
}

type NewEventType struct {
	Name string `json:"name" bson:"name"`
}

type NewFacility struct {
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
	Type string `json:"type" bson:"type"`
}

type NewFacilityHistory struct {
	ID         *string   `json:"id" bson:"_id"`
	FacilityID string    `json:"facilityId" bson:"facilityId"`
	BorrowDate time.Time `json:"borrowDate" bson:"borrowDate"`
	ReturnDate time.Time `json:"returnDate" bson:"returnDate"`
	EventID    *string   `json:"eventId" bson:"eventId"`
}

type NewParticipant struct {
	EventID              string    `json:"eventId" bson:"eventId"`
	Email                string    `json:"email" bson:"email"`
	Name                 string    `json:"name" bson:"name"`
	Academic             string    `json:"academic" bson:"academic"`
	School               string    `json:"school" bson:"school"`
	Major                string    `json:"major" bson:"major"`
	Phone                string    `json:"phone" bson:"phone"`
	Dob                  time.Time `json:"dob" bson:"dob"`
	ExpectedGraduateDate time.Time `json:"expectedGraduateDate" bson:"expectedGraduateDate"`
}

type NewTask struct {
	ID        *string   `json:"id" bson:"_id"`
	EventID   *string   `json:"eventId" bson:"eventId"`
	Name      string    `json:"name" bson:"name"`
	UserID    string    `json:"userId" bson:"userId"`
	Type      string    `json:"type" bson:"type"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
}

type NewUser struct {
	Email           string   `json:"email" bson:"email"`
	Password        string   `json:"password" bson:"password"`
	ConfirmPassword string   `json:"confirmPassword" bson:"confirmPassword"`
	Roles           []string `json:"roles" bson:"roles"`
}

type Participant struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt            time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsValid              bool               `json:"isValid" bson:"isValid"`
	IsAttended           bool               `json:"isAttended" bson:"isAttended"`
	Event                *Event             `json:"event" bson:"event"`
	Email                string             `json:"email" bson:"email"`
	Name                 string             `json:"name" bson:"name"`
	Academic             string             `json:"academic" bson:"academic"`
	School               string             `json:"school" bson:"school"`
	Major                string             `json:"major" bson:"major"`
	Phone                string             `json:"phone" bson:"phone"`
	Dob                  time.Time          `json:"dob" bson:"dob"`
	ExpectedGraduateDate time.Time          `json:"expectedGraduateDate" bson:"expectedGraduateDate"`
}

type Task struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	Event     *Event             `json:"event" bson:"event"`
	Name      string             `json:"name" bson:"name"`
	User      *User              `json:"user" bson:"user"`
	Type      string             `json:"type" bson:"type"`
	StartDate time.Time          `json:"startDate" bson:"startDate"`
	EndDate   time.Time          `json:"endDate" bson:"endDate"`
}

type UpdateEvent struct {
	Tags                  []string               `json:"tags" bson:"tags"`
	Tasks                 []*NewTask             `json:"tasks" bson:"tasks"`
	FacilityHistories     []*NewFacilityHistory  `json:"facilityHistories" bson:"facilityHistories"`
	Name                  string                 `json:"name" bson:"name"`
	Language              string                 `json:"language" bson:"language"`
	EventTypeID           string                 `json:"eventTypeId" bson:"eventTypeId"`
	Mode                  string                 `json:"mode" bson:"mode"`
	Location              string                 `json:"location" bson:"location"`
	Accommodation         string                 `json:"accommodation" bson:"accommodation"`
	RegistrationCloseDate time.Time              `json:"registrationCloseDate" bson:"registrationCloseDate"`
	StartDate             time.Time              `json:"startDate" bson:"startDate"`
	EndDate               time.Time              `json:"endDate" bson:"endDate"`
	MaxParticipants       int                    `json:"maxParticipants" bson:"maxParticipants"`
	Description           string                 `json:"description" bson:"description"`
	OwnerID               string                 `json:"ownerId" bson:"ownerId"`
	Budget                float64                `json:"budget" bson:"budget"`
	Image                 string                 `json:"image" bson:"image"`
	IsApproved            bool                   `json:"isApproved" bson:"isApproved"`
	Reviewer              *string                `json:"reviewer" bson:"reviewer"`
	IsFinished            bool                   `json:"isFinished" bson:"isFinished"`
	IsDeleted             bool                   `json:"isDeleted" bson:"isDeleted"`
	CustomizeFields       []*InputCustomizeField `json:"customizeFields" bson:"customizeFields"`
}

type UpdateEventType struct {
	Name      string `json:"name" bson:"name"`
	IsDeleted bool   `json:"isDeleted" bson:"isDeleted"`
}

type UpdateFacility struct {
	Name      string `json:"name" bson:"name"`
	Code      string `json:"code" bson:"code"`
	Type      string `json:"type" bson:"type"`
	Status    bool   `json:"status" bson:"status"`
	IsDeleted bool   `json:"isDeleted" bson:"isDeleted"`
}

type UpdateFacilityHistory struct {
	FacilityID string    `json:"facilityId" bson:"facilityId"`
	BorrowDate time.Time `json:"borrowDate" bson:"borrowDate"`
	ReturnDate time.Time `json:"returnDate" bson:"returnDate"`
	EventID    string    `json:"eventId" bson:"eventId"`
}

type UpdateParticipant struct {
	EventID              string    `json:"eventId" bson:"eventId"`
	Email                string    `json:"email" bson:"email"`
	Name                 string    `json:"name" bson:"name"`
	Academic             string    `json:"academic" bson:"academic"`
	School               string    `json:"school" bson:"school"`
	Major                string    `json:"major" bson:"major"`
	Phone                string    `json:"phone" bson:"phone"`
	Dob                  time.Time `json:"dob" bson:"dob"`
	ExpectedGraduateDate time.Time `json:"expectedGraduateDate" bson:"expectedGraduateDate"`
	IsValid              bool      `json:"isValid" bson:"isValid"`
	IsAttended           bool      `json:"isAttended" bson:"isAttended"`
}

type UpdateTask struct {
	EventID   *string   `json:"eventId" bson:"eventId"`
	Name      string    `json:"name" bson:"name"`
	UserID    string    `json:"userId" bson:"userId"`
	Type      string    `json:"type" bson:"type"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
}

type UpdateUser struct {
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"password" bson:"password"`
	Roles    []string `json:"roles" bson:"roles"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Roles     []string           `json:"roles" bson:"roles"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
