package controllers

import (
	"fmt"
	"go-fiber/config"
	"go-fiber/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}

	auth := models.Auth{
		Username: data["username"],
		Password: string(hashedPassword),
	}

	// Cek apakah create berhasil
	result := config.DB.Create(&auth)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   result.Error.Error(),
		})
	}

	user := models.User{
		Name:   data["name"],  // pastikan dikirim di body request
		Email:  data["email"], // pastikan dikirim juga
		AuthID: uint(auth.ID),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create user profile",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "User registered successfully",
		"auth_id":  auth.ID,
		"user_id":  user.ID,
		"username": auth.Username,
		"name":     user.Name,
		"email":    user.Email,
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	var auth models.Auth

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Bad Request",
		})
	}

	config.DB.Where("username = ?", data["username"]).First(&auth)

	if auth.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(data["password"]))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  auth.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(config.ENV.JWTSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"token": token})
}

func Logout(c *fiber.Ctx) error {
	// hapus cookie token
	token := c.Cookies("token")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No token found",
		})
	}

	fmt.Println("Token to blacklist:", token)

	err := config.RDB.Set(config.Ctx, "blacklist:"+token, "1", time.Hour*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to blacklist token",
		})
	}

	c.ClearCookie("token")

	return c.JSON(fiber.Map{
		"message": "Logout berhasil dan token di-blacklist",
	})
}
