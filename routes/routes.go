package routes

import (
	"go-fiber/controllers"
	"go-fiber/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(route fiber.Router) {
	// Add logging middleware for debugging
	route.Use(func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	// Public routes (no authentication required)
	route.Post("/register", controllers.Register)
	route.Post("/login", controllers.Login)

	// Handle GET requests to register/login (wrong method)
	route.Get("/register", func(c *fiber.Ctx) error {
		return c.Status(405).JSON(fiber.Map{
			"message": "Method not allowed. Use POST for registration.",
			"method":  c.Method(),
			"path":    c.Path(),
		})
	})

	route.Get("/login", func(c *fiber.Ctx) error {
		return c.Status(405).JSON(fiber.Map{
			"message": "Method not allowed. Use POST for login.",
			"method":  c.Method(),
			"path":    c.Path(),
		})
	})

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
