package routes

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/khanhvtn/netevent-go/controllers"
// )

// /* SetupServerRoutes: setup all routes for the server */
// func SetupServerRoutes(app *fiber.App) {
// 	api := app.Group("/api")
// 	setUserRoutes(api)
// }

// /* User Routes */
// func setUserRoutes(route fiber.Router) {
// 	userRoute := route.Group("/users")
// 	userRoute.Post("/login", controllers.Login)
// 	userRoute.Delete("/logout", controllers.Logout)
// 	userRoute.Get("/checkUser", controllers.CheckUser)

// 	userRoute.Get("/", controllers.GetAllUsers)
// 	userRoute.Get("/:id", controllers.GetAnUser)
// 	userRoute.Post("/", controllers.CreateAnUser)
// 	userRoute.Patch("/:id", controllers.UpdateAnUser)
// 	userRoute.Delete("/:id", controllers.DeleteAnUser)

// }
