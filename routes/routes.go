package routes

import (
	"go-fiber/controllers"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(route fiber.Router) {
	// Public routes (no authentication required)
	route.Post("/register", controllers.Register)
	route.Post("/login", controllers.Login)

	// Protected routes (authentication required)
	route.Use(middleware.Protected())

	route.Post("/logout", controllers.Logout)
	route.Get("/profile", controllers.GetProfile)
	route.Get("/", controllers.GetUsers)
	route.Get("/:id", controllers.GetUserById)
	route.Post("/", controllers.CreateUser)
	route.Put("/:id", controllers.UpdateUser)
	route.Delete("/:id", controllers.DeleteUser)
}
