package main

import (
	"go-fiber/config"
	"go-fiber/controllers"
	"go-fiber/models"
	"go-fiber/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Load .env dulu ke struct
	config.LoadEnv()
	config.LoadOPAConfig()

	// 2. Connect ke DB pakai isi dari struct
	config.ConnectDB()

	// 3. Auto migrate models
	config.DB.AutoMigrate(
		&models.User{},
		&models.Conversation{},
		&models.Message{},
		&models.ChatHistory{},
	)

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50MB limit for file uploads
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Serve static files (uploaded images)
	app.Static("/uploads", "./uploads")

	// 4. Routing
	api := app.Group("/api")
	
	// Health check (no auth required)
	app.Get("/health", controllers.HealthCheck)
	
	// User routes (existing)
	user := api.Group("/user")
	routes.UserRoute(user)
	
	// Chat routes (new)
	chat := api.Group("/chat")
	routes.ChatRoute(chat)

	// 5. Listen pakai port dari .env
	app.Listen(":" + config.ENV.AppPort)
}
