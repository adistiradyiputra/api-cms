package controllers

import (
	"go-fiber/config"
	"go-fiber/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	config.DB.Find(&users)
	return c.JSON(users)
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	result := config.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	return c.JSON(user)
}

func GetProfile(c *fiber.Ctx) error {
	// Ambil token JWT dari middleware
	userToken := c.Locals("user")
	if userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse klaim dari token JWT
	claims := userToken.(*jwt.Token).Claims.(jwt.MapClaims)
	authID := uint(claims["id"].(float64))

	// Cari user berdasarkan auth_id
	var user models.User
	if err := config.DB.Where("auth_id = ?", authID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Cari auth berdasarkan authID untuk mendapatkan username
	var auth models.Auth
	if err := config.DB.Where("id = ?", authID).First(&auth).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Auth not found",
		})
	}

	// Buat response yang menggabungkan data user dan auth
	response := fiber.Map{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"auth_id":  user.AuthID,
		"username": auth.Username,
	}

	return c.JSON(response)
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	config.DB.Create(&user)
	return c.JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}
	config.DB.Save(&user)
	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	result := config.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}
