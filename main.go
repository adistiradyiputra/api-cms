package main

import (
	"go-fiber/config"
	"go-fiber/models"
	"go-fiber/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	log.Println("Starting application...")

	// 1. Load .env dulu ke struct
	config.LoadEnv()
	log.Println("Environment loaded")

	// Connect Redis (optional)
	config.ConnectRedis()
	log.Println("Redis connection attempted")

	// 2. Connect ke DB pakai isi dari struct (optional untuk testing)
	config.ConnectDB()

	// 3. Auto migrate user model (only if DB is connected)
	if config.DB != nil {
		// Force drop and recreate tables for testing
		log.Println("Dropping existing tables...")
		config.DB.Migrator().DropTable(&models.User{}, &models.Auth{})

		log.Println("Creating tables with new schema...")
		config.DB.AutoMigrate(&models.User{}, &models.Auth{})
		log.Println("Database migration completed")
	} else {
		log.Println("Warning: Database not connected, skipping migration")
		log.Println("Application will start without database connection")
	}

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

	// 5. Listen pakai port dari Railway atau .env atau default 8080
	port := os.Getenv("PORT") // Railway menggunakan PORT
	if port == "" {
		port = config.ENV.AppPort
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
