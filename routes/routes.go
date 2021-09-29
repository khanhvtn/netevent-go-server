package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/khanhvtn/netevent-go/controllers"
)

/* SetupServerRoutes: setup all routes for the server */
func SetupServerRoutes(app *gin.Engine) {
	api := app.Group("/api")
	setUserRoutes(api)
}

/* User Routes */
func setUserRoutes(api *gin.RouterGroup) {
	eventRoute := api.Group("/event")
	eventRoute.GET("/eventStatistic/:id", controllers.GetEventStatistic)
}
