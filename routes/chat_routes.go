package routes

import (
	"go-fiber/controllers"
	"go-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func ChatRoute(route fiber.Router) {
	// Apply OPA.co.id authentication middleware to all chat routes
	route.Use(middleware.ValidateOPAToken)

	// Chat endpoints
	route.Post("/send", controllers.SendMessage)           // Send message and get response
	route.Post("/stream", controllers.StreamMessage)       // Stream chat response
	route.Post("/save", controllers.SaveChat)              // Save chat data
	route.Delete("/conversation", controllers.DeleteConversation) // Delete conversation
}