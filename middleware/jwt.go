package middleware

import (
	"go-fiber/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Protected() fiber.Handler {
	// Check if JWT secret is set
	if config.ENV.JWTSecret == "" {
		// If no JWT secret, allow all requests (for testing)
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return jwtware.New(jwtware.Config{
		SigningKey:  []byte(config.ENV.JWTSecret),
		TokenLookup: "cookie:token",

		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Cookies("token")
			if token == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "No token provided",
				})
			}

			// Skip Redis check if Redis is not connected
			if config.RDB != nil {
				_, err := config.RDB.Get(config.Ctx, "blacklist:"+token).Result()
				if err == nil {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"message": "Token sudah tidak berlaku (blacklisted)",
					})
				}
			}

			return c.Next()
		},

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		},
	})
}
