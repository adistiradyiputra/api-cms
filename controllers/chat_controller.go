package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-fiber/config"
	"go-fiber/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ChatRequest represents the incoming chat request
type ChatRequest struct {
	Content         string `json:"content" form:"content"`
	ConversationID  string `json:"conversation_id" form:"conversation_id"`
	Persona         string `json:"persona" form:"persona"`
	ResponseMode    string `json:"response_mode" form:"response_mode"`
	Reranker        string `json:"reranker" form:"reranker"`
	ModelName       string `json:"model_name" form:"model_name"`
	StreamMessage   string `json:"stream_message" form:"stream_message"`
}

// APIResponse represents the response from the chat API
type APIResponse struct {
	Status string `json:"status"`
	Data   struct {
		Message          string `json:"message"`
		ConversationID   string `json:"conversation_id"`
		Content          string `json:"content"`
		Role             string `json:"role"`
	} `json:"data"`
}

// TokenResponse represents streaming token response
type TokenResponse struct {
	Token string `json:"token"`
}

// SendMessage handles sending messages to the chat API
func SendMessage(c *fiber.Ctx) error {
	// Get user data from context (set by middleware)
	userID := c.Locals("user_id").(uint)
	sessionID := c.Locals("session_id").(string)

	// Parse request
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
			"detail": err.Error(),
		})
	}

	// Check for file upload
	file, err := c.FormFile("image")
	var imageFileName, imageFileURL, imageFileType string
	var imageData []byte

	if err == nil && file != nil {
		// Handle image upload
		imageFileName = file.Filename
		imageFileType = file.Header.Get("Content-Type")
		
		// Create uploads directory if it doesn't exist
		uploadDir := "./uploads/images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create upload directory",
			})
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, newFileName)

		// Save file
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save image",
			})
		}

		imageFileURL = fmt.Sprintf("/uploads/images/%s", newFileName)
		
		// Read file data for API request
		imageData, err = os.ReadFile(filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read image data",
			})
		}

		// Set default content if empty
		if req.Content == "" {
			req.Content = "[Image]"
		}
	}

	// Validate content
	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content is required",
		})
	}

	// Set default values
	if req.Persona == "" {
		req.Persona = "Normal"
	}
	if req.ResponseMode == "" {
		req.ResponseMode = "short"
	}
	if req.Reranker == "" {
		req.Reranker = "false"
	}
	if req.ModelName == "" {
		req.ModelName = "llama-4"
	}

	// Handle conversation ID
	apiConversationID := "0"
	if req.ConversationID != "" && req.ConversationID != "0" {
		// Check if conversation exists
		var conversation models.Conversation
		result := config.DB.Where("conversation_id = ?", req.ConversationID).First(&conversation)
		if result.Error == nil && conversation.APIConversationID != "" {
			apiConversationID = conversation.APIConversationID
		}
	}

	// Check conversation limit
	var messageCount int64
	config.DB.Model(&models.Message{}).Where("conversation_id = ?", req.ConversationID).Count(&messageCount)
	if messageCount > 50 {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":  "Conversation limit reached",
			"detail": "Silakan buat conversation baru",
		})
	}

	// Prepare multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add form fields
	fields := map[string]string{
		"content":         req.Content,
		"conversation_id": apiConversationID,
		"persona":         req.Persona,
		"response_mode":   req.ResponseMode,
		"reranker":        req.Reranker,
		"model_name":      req.ModelName,
	}

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to prepare request",
			})
		}
	}

	// Add image if present
	if len(imageData) > 0 {
		part, err := writer.CreateFormFile("image", imageFileName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to prepare image data",
			})
		}
		part.Write(imageData)
	}

	writer.Close()

	// Make API request
	apiURL := "http://8.215.40.27:5022/chat-stream/"
	apiReq, err := http.NewRequest("POST", apiURL, &requestBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create API request",
		})
	}

	apiReq.Header.Set("Content-Type", writer.FormDataContentType())
	apiReq.Header.Set("Accept", "text/event-stream")
	apiReq.Header.Set("x-api-key", "WzkprPudy3zN4UyzsMoiRVOh2B3uCntaDfHyAbww03QRfheNQALmIOlVIjGOGdHu")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(apiReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to chat API",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":  "API Error",
			"detail": "Gagal terhubung ke API streaming",
		})
	}

	// Process streaming response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read API response",
		})
	}

	// Parse SSE response
	lines := strings.Split(string(body), "\n")
	var tokenContent string
	var finalResponse APIResponse

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if strings.Contains(data, `"token"`) {
				var tokenResp TokenResponse
				if json.Unmarshal([]byte(data), &tokenResp) == nil {
					tokenContent += tokenResp.Token
				}
			} else if strings.Contains(data, `"status"`) {
				json.Unmarshal([]byte(data), &finalResponse)
			}
		}
	}

	// Save to database
	if err := saveChatToDatabase(c, req, userID, sessionID, tokenContent, finalResponse, imageFileName, imageFileURL, imageFileType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save chat",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Pesan berhasil dikirim dan disimpan.",
		"data": fiber.Map{
			"msg": fiber.Map{
				"content": tokenContent,
				"role":    "assistant",
			},
			"conversation_id":      req.ConversationID,
			"api_conversation_id":  finalResponse.Data.ConversationID,
			"session_id":          sessionID,
		},
	})
}

// StreamMessage handles streaming chat responses
func StreamMessage(c *fiber.Ctx) error {
	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	// Get user data from context
	userID := c.Locals("user_id").(uint)
	sessionID := c.Locals("session_id").(string)

	// Parse request
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		sendSSEError(c, "Invalid request format")
		return nil
	}

	// Validate content
	if req.Content == "" {
		sendSSEError(c, "Content is required")
		return nil
	}

	// Handle file upload (similar to SendMessage)
	file, err := c.FormFile("image")
	var imageFileName, imageFileURL, imageFileType string
	var imageData []byte

	if err == nil && file != nil {
		imageFileName = file.Filename
		imageFileType = file.Header.Get("Content-Type")
		
		uploadDir := "./uploads/images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			sendSSEError(c, "Failed to create upload directory")
			return nil
		}

		ext := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, newFileName)

		if err := c.SaveFile(file, filePath); err != nil {
			sendSSEError(c, "Failed to save image")
			return nil
		}

		imageFileURL = fmt.Sprintf("/uploads/images/%s", newFileName)
		
		imageData, err = os.ReadFile(filePath)
		if err != nil {
			sendSSEError(c, "Failed to read image data")
			return nil
		}

		if req.Content == "" {
			req.Content = "[Image]"
		}
	}

	// Set default values
	if req.Persona == "" {
		req.Persona = "Normal"
	}
	if req.ResponseMode == "" {
		req.ResponseMode = "short"
	}
	if req.Reranker == "" {
		req.Reranker = "false"
	}
	if req.ModelName == "" {
		req.ModelName = "llama-4"
	}

	// Handle conversation ID
	apiConversationID := "0"
	if req.ConversationID != "" && req.ConversationID != "0" {
		var conversation models.Conversation
		result := config.DB.Where("conversation_id = ?", req.ConversationID).First(&conversation)
		if result.Error == nil && conversation.APIConversationID != "" {
			apiConversationID = conversation.APIConversationID
		}
	}

	// Prepare multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	fields := map[string]string{
		"content":         req.Content,
		"conversation_id": apiConversationID,
		"persona":         req.Persona,
		"response_mode":   req.ResponseMode,
		"reranker":        req.Reranker,
		"model_name":      req.ModelName,
	}

	for key, value := range fields {
		writer.WriteField(key, value)
	}

	if len(imageData) > 0 {
		part, _ := writer.CreateFormFile("image", imageFileName)
		part.Write(imageData)
	}

	writer.Close()

	// Make API request
	apiURL := "http://8.215.40.27:5022/chat-stream/"
	apiReq, err := http.NewRequest("POST", apiURL, &requestBody)
	if err != nil {
		sendSSEError(c, "Failed to create API request")
		return nil
	}

	apiReq.Header.Set("Content-Type", writer.FormDataContentType())
	apiReq.Header.Set("Accept", "text/event-stream")
	apiReq.Header.Set("x-api-key", "WzkprPudy3zN4UyzsMoiRVOh2B3uCntaDfHyAbww03QRfheNQALmIOlVIjGOGdHu")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(apiReq)
	if err != nil {
		sendSSEError(c, "Failed to connect to chat API")
		return nil
	}
	defer resp.Body.Close()

	// Stream response back to client
	buffer := make([]byte, 1024)
	var tokenContent string
	var finalResponse APIResponse

	for {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			break
		}

		if n > 0 {
			chunk := string(buffer[:n])
			lines := strings.Split(chunk, "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					// Forward the SSE event to client
					c.Write([]byte(line + "\n"))
					c.Context().Response.Flush()

					// Parse for final response
					if strings.HasPrefix(line, "data: ") {
						data := strings.TrimPrefix(line, "data: ")
						if strings.Contains(data, `"token"`) {
							var tokenResp TokenResponse
							if json.Unmarshal([]byte(data), &tokenResp) == nil {
								tokenContent += tokenResp.Token
							}
						} else if strings.Contains(data, `"status"`) {
							json.Unmarshal([]byte(data), &finalResponse)
						}
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	// Save to database after streaming
	if finalResponse.Status == "success" {
		if err := saveChatToDatabase(c, req, userID, sessionID, tokenContent, finalResponse, imageFileName, imageFileURL, imageFileType); err != nil {
			sendSSEError(c, "Failed to save chat")
			return nil
		}

		// Send save success event
		saveData := fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"conversation_id":     req.ConversationID,
				"api_conversation_id": finalResponse.Data.ConversationID,
			},
		}
		saveJSON, _ := json.Marshal(saveData)
		c.Write([]byte(fmt.Sprintf("event: save_success\ndata: %s\n\n", string(saveJSON))))
		c.Context().Response.Flush()
	}

	return nil
}

// DeleteConversation handles conversation deletion
func DeleteConversation(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req struct {
		ConversationID string `json:"conversation_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.ConversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Conversation ID is required",
		})
	}

	// Check if conversation exists and belongs to user
	var conversation models.Conversation
	result := config.DB.Where("conversation_id = ? AND user_id = ?", req.ConversationID, userID).First(&conversation)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":  "Not Found",
			"detail": "Conversation not found or you do not have permission to delete it",
		})
	}

	// Soft delete conversation
	if err := config.DB.Delete(&conversation).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete conversation",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Conversation successfully deleted",
		"data": fiber.Map{
			"conversation_id": req.ConversationID,
		},
	})
}

// SaveChat handles saving chat data (for when streaming is handled client-side)
func SaveChat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	sessionID := c.Locals("session_id").(string)

	var req struct {
		Content         string `json:"content"`
		StreamMessage   string `json:"stream_message"`
		ConversationID  string `json:"conversation_id"`
		APIConversationID string `json:"api_conversation_id"`
		ImageName       string `json:"image_name"`
		ImageURL        string `json:"image_url"`
		ImageType       string `json:"image_type"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Content == "" || req.StreamMessage == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content and stream_message are required",
		})
	}

	// Create chat request for database save
	chatReq := ChatRequest{
		Content: req.Content,
		ConversationID: req.ConversationID,
	}

	// Create mock API response
	apiResponse := APIResponse{
		Status: "success",
		Data: struct {
			Message          string `json:"message"`
			ConversationID   string `json:"conversation_id"`
			Content          string `json:"content"`
			Role             string `json:"role"`
		}{
			Message: req.StreamMessage,
			ConversationID: req.APIConversationID,
			Content: req.StreamMessage,
			Role: "assistant",
		},
	}

	// Save to database
	if err := saveChatToDatabase(c, chatReq, userID, sessionID, req.StreamMessage, apiResponse, req.ImageName, req.ImageURL, req.ImageType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save chat",
			"detail": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Pesan berhasil disimpan.",
		"data": fiber.Map{
			"msg": fiber.Map{
				"content": req.StreamMessage,
				"role":    "assistant",
			},
			"conversation_id":     req.ConversationID,
			"api_conversation_id": req.APIConversationID,
			"session_id":         sessionID,
		},
	})
}

// Helper functions

func saveChatToDatabase(c *fiber.Ctx, req ChatRequest, userID uint, sessionID, tokenContent string, finalResponse APIResponse, imageFileName, imageFileURL, imageFileType string) error {
	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Handle conversation
	var conversation models.Conversation
	var conversationID string

	if req.ConversationID != "" && req.ConversationID != "0" {
		// Check if conversation exists
		result := tx.Where("conversation_id = ?", req.ConversationID).First(&conversation)
		if result.Error == nil {
			conversationID = conversation.ConversationID
			// Update timestamp
			tx.Model(&conversation).Update("updated_at", time.Now())
		}
	}

	if conversationID == "" {
		// Create new conversation
		conversationID = generateConversationID()
		conversation = models.Conversation{
			ConversationID:    conversationID,
			APIConversationID: finalResponse.Data.ConversationID,
			SessionID:         sessionID,
			UserID:            userID,
			Title:             truncateString(req.Content, 100),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		if err := tx.Create(&conversation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create message metadata
	messageMetadata := models.JSONMap{
		"has_image":  imageFileName != "",
		"image_name": imageFileName,
		"image_url":  imageFileURL,
		"image_type": imageFileType,
	}

	// Save user message
	userMessageID := generateMessageID()
	userMessage := models.Message{
		MessageID:       userMessageID,
		ConversationID:  conversationID,
		ParentMessageID: nil,
		Role:            "user",
		Content:         req.Content,
		MessageMetadata: messageMetadata,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(&userMessage).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Save assistant message
	assistantMessageID := generateMessageID()
	assistantMessage := models.Message{
		MessageID:       assistantMessageID,
		ConversationID:  conversationID,
		ParentMessageID: &userMessageID,
		Role:            "assistant",
		Content:         tokenContent,
		MessageMetadata: models.JSONMap{
			"api_response": finalResponse,
			"timestamp":    time.Now().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tx.Create(&assistantMessage).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Save chat history
	chatHistory := models.ChatHistory{
		ChatID:                generateChatID(),
		ConversationID:        conversationID,
		ConversationSessionID: sessionID,
		MessageUser:           req.Content,
		MessageAssistant:      tokenContent,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		FileName:              &imageFileName,
		FileURL:               &imageFileURL,
		FileType:              &imageFileType,
	}

	if err := tx.Create(&chatHistory).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func sendSSEError(c *fiber.Ctx, message string) {
	errorData := fiber.Map{"error": message}
	errorJSON, _ := json.Marshal(errorData)
	c.Write([]byte(fmt.Sprintf("event: error\ndata: %s\n\n", string(errorJSON))))
	c.Context().Response.Flush()
}

func generateConversationID() string {
	return "conv_" + uuid.New().String()
}

func generateMessageID() string {
	return "msg_" + uuid.New().String()
}

func generateChatID() string {
	return "chat_" + uuid.New().String()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}