package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mattchw/go-onboard/controllers"
	middleware "github.com/mattchw/go-onboard/middlewares"
)

func UserRoute(app *fiber.App) {
	app.Get("/users", controllers.GetUsers)
	app.Post("/users", middleware.AuthReq(), controllers.CreateUser)
	app.Get("/users/:userId", controllers.GetUser)
	app.Patch("/users/:userId", controllers.UpdateUser)
	app.Delete("/users/:userId", controllers.DeleteUser)
}
