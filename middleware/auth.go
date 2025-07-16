package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go-fiber/config"
	"github.com/gofiber/fiber/v2"
)

type OPATokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		UserID    uint   `json:"user_id"`
		SessionID string `json:"session_id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
	} `json:"data"`
}

// ValidateOPAToken validates the access token from OPA.co.id
func ValidateOPAToken(c *fiber.Ctx) error {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"detail":  "Authorization header is required",
			"message": "Silakan login kembali untuk melanjutkan",
		})
	}

	// Check if it's a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"detail":  "Invalid authorization format",
			"message": "Silakan login kembali untuk melanjutkan",
		})
	}

	// Extract token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate token with OPA.co.id
	userData, err := validateTokenWithOPA(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"detail":  "Invalid or expired token",
			"message": "Silakan login kembali untuk melanjutkan",
		})
	}

	// Store user data in context for later use
	c.Locals("user_id", userData.Data.UserID)
	c.Locals("session_id", userData.Data.SessionID)
	c.Locals("user_email", userData.Data.Email)
	c.Locals("user_name", userData.Data.Name)

	return c.Next()
}

// validateTokenWithOPA sends request to OPA.co.id to validate the token
func validateTokenWithOPA(token string) (*OPATokenResponse, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create request to OPA.co.id validation endpoint
	req, err := http.NewRequest("POST", config.OPA.APIBaseURL+"/validate-token", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResponse OPATokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if token is valid
	if tokenResponse.Status != "success" {
		return nil, fmt.Errorf("token validation failed: %s", tokenResponse.Message)
	}

	return &tokenResponse, nil
}