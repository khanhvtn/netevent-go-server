package controllers

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khanhvtn/netevent-go/models"
	"github.com/khanhvtn/netevent-go/services"
	"github.com/khanhvtn/netevent-go/utilities"
	"github.com/sarulabs/di"
	"go.mongodb.org/mongo-driver/bson"
)

func handleError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
}
func GetEventStatistic(c *gin.Context) {
	idEvent := c.Param("id")
	container := c.MustGet("container").(di.Container)
	eventService := container.Get(services.EventServiceName).(*services.EventService)
	eventTypeService := container.Get(services.EventTypeServiceName).(*services.EventTypeService)
	userService := container.Get(services.UserServiceName).(*services.UserService)
	taskService := container.Get(services.TaskServiceName).(*services.TaskService)
	facilityHistoryService := container.Get(services.FacilityHistoryServiceName).(*services.FacilityHistoryService)
	facilityService := container.Get(services.FacilityServiceName).(*services.FacilityService)
	objectID, err := utilities.ConvertStringIdToObjectID(idEvent)
	if err != nil {
		handleError(c, err)
	}

	event, err := eventService.GetOne(bson.M{"_id": objectID})
	if err != nil {
		handleError(c, err)
	}
	file, err := os.CreateTemp("", "random")
	if err != nil {
		handleError(c, err)
	}

	defer file.Close()
	defer os.Remove(file.Name())
	//generate csv file based on event information
	if _, err := file.WriteString("Name;CreatedAt;UpdatedAt;Tags;IsApproved;Reviewer;IsFinished;Tasks;FacilityHistories;Language;EventType;Mode;Location;Accommodation;RegistrationCloseDate;StartDate;EndDate;MaxParticipants;Description;Owner;Budget;Image;IsDeleted;CustomizeFields\n"); err != nil {
		handleError(c, err)
	}

	eventType, err := eventTypeService.GetOne(bson.M{"_id": event.EventType})
	if err != nil {
		handleError(c, err)
	}

	owner, err := userService.GetOne(bson.M{"_id": event.Owner})
	if err != nil {
		handleError(c, err)
	}
	reviewer, err := userService.GetOne(bson.M{"_id": event.Reviewer})
	if err != nil {
		handleError(c, err)
	}

	//generate task
	tasksChan := make(chan string)
	go func(tasksChan chan string, event *models.Event) {
		tasks := make([]string, 0)
		for _, id := range event.Tasks {
			task, err := taskService.GetOne(bson.M{"_id": id})
			if err != nil {
				handleError(c, err)
			}
			tasks = append(tasks, task.Name)
		}
		tasksChan <- strings.Join(tasks, ",")
	}(tasksChan, event)

	//generate task
	facilityHistoriesChan := make(chan string)
	go func(facilityHistoriesChan chan string, event *models.Event) {
		facilities := make([]string, 0)
		for _, id := range event.FacilityHistories {
			facilityHistory, err := facilityHistoryService.GetOne(bson.M{"_id": id})
			if err != nil {
				handleError(c, err)
			}
			facility, err := facilityService.GetOne(bson.M{"_id": facilityHistory.Facility})
			if err != nil {
				handleError(c, err)
			}
			facilities = append(facilities, facility.Name)
		}

		facilityHistoriesChan <- strings.Join(facilities, ",")
	}(facilityHistoriesChan, event)

	//generate event customize fields
	customizeFieldsChan := make(chan string)
	go func(customizeFieldsChan chan string, event *models.Event) {
		customizeFields := make([]string, 0)
		for _, field := range event.CustomizeFields {
			customizeFields = append(customizeFields, field.Name)
		}
		strings.Join(customizeFields, ",")
		customizeFieldsChan <- strings.Join(customizeFields, ",")
	}(customizeFieldsChan, event)

	if _, err := file.WriteString(fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%f;%v;%v;%v\n", event.Name, event.CreatedAt, event.UpdatedAt, event.Tags, event.IsApproved, reviewer.Email, event.IsFinished, <-tasksChan, <-facilityHistoriesChan, event.Language, eventType.Name, event.Mode, event.Location, event.Accommodation, event.RegistrationCloseDate, event.StartDate, event.EndDate, event.MaxParticipants, event.Description, owner.Email, event.Budget, event.Image, event.IsDeleted, <-customizeFieldsChan)); err != nil {
		handleError(c, err)
	}

	if err != nil {
		handleError(c, err)
	}
	c.FileAttachment(file.Name(), "event.csv")
}
