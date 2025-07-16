package models

import (
	"time"

	"gorm.io/gorm"
)

type ChatHistory struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	ChatID                string         `json:"chat_id" gorm:"uniqueIndex;not null"`
	ConversationID        string         `json:"conversation_id" gorm:"index;not null"`
	ConversationSessionID string         `json:"conversation_session_id" gorm:"index;not null"`
	MessageUser           string         `json:"message_user" gorm:"type:text"`
	MessageAssistant      string         `json:"message_assistant" gorm:"type:text"`
	PreviousChatID        *string        `json:"previous_chat_id" gorm:"index"`
	FileName              *string        `json:"file_name" gorm:"size:255"`
	FileURL               *string        `json:"file_url" gorm:"size:500"`
	FileType              *string        `json:"file_type" gorm:"size:100"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// Relationships
	Conversation Conversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationID;references:ConversationID"`
}