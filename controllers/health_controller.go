package controllers

import (
	"go-fiber/config"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck returns the health status of the API
func HealthCheck(c *fiber.Ctx) error {
	// Check database connection
	dbStatus := "healthy"
	if err := config.DB.Raw("SELECT 1").Error; err != nil {
		dbStatus = "unhealthy"
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"message": "Chat API with OPA.co.id Integration is running",
		"data": fiber.Map{
			"database": dbStatus,
			"version":  "1.0.0",
			"features": []string{
				"OPA.co.id Authentication",
				"Chat Functionality",
				"File Upload",
				"Streaming Support",
				"Database Storage",
			},
		},
	})
}