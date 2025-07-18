package routes

import (
	"go-fiber/controllers"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(route fiber.Router) {
	// Auth
	route.Post("/register", controllers.Register)
	route.Post("/login", controllers.Login)
	route.Post("/logout", controllers.Logout)

	route.Use(middleware.Protected())

	route.Get("/profile", controllers.GetProfile)

	route.Get("/", controllers.GetUsers)
	route.Get("/:id", controllers.GetUserById)
	route.Post("/", controllers.CreateUser)
	route.Put("/:id", controllers.UpdateUser)
	route.Delete("/:id", controllers.DeleteUser)
}
