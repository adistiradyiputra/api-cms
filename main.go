package main

import (
	"go-fiber/config"
	"go-fiber/models"
	"go-fiber/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// 2. Connect ke DB pakai isi dari struct
	config.ConnectDB()

	// 3. Auto migrate user model
	config.DB.AutoMigrate(&models.User{})

	app := fiber.New()

	// 4. Routing
	api := app.Group("/api")
	user := api.Group("/user")
	routes.UserRoute(user)

	// 5. Listen pakai port dari .env
	app.Listen(":" + config.ENV.AppPort)
}
