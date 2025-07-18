package main

import (
	"go-fiber/config"
	"go-fiber/models"
	"go-fiber/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Load .env dulu ke struct
	config.LoadEnv()
	config.ConnectRedis()

	// 2. Connect ke DB pakai isi dari struct
	config.ConnectDB()

	// 3. Auto migrate user model
	config.DB.AutoMigrate(&models.User{}, &models.Auth{})

	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // Vue.js dev server URL
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// 4. Routing
	// Health check route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	api := app.Group("/api")
	user := api.Group("/user")
	routes.UserRoute(user)

	// 5. Listen pakai port dari .env
	app.Listen(":" + config.ENV.AppPort)
}
