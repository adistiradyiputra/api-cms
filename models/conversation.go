package models

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	ConversationID      string         `json:"conversation_id" gorm:"uniqueIndex;not null"`
	APIConversationID   string         `json:"api_conversation_id" gorm:"index"`
	SessionID           string         `json:"session_id" gorm:"index;not null"`
	UserID              uint           `json:"user_id" gorm:"index;not null"`
	Title               string         `json:"title" gorm:"size:255"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// Relationships
	Messages     []Message      `json:"messages,omitempty" gorm:"foreignKey:ConversationID;references:ConversationID"`
	ChatHistories []ChatHistory `json:"chat_histories,omitempty" gorm:"foreignKey:ConversationID;references:ConversationID"`
}