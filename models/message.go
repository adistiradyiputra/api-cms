package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSONMap is a custom type to handle JSON metadata
type JSONMap map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &j)
}

type Message struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	MessageID       string         `json:"message_id" gorm:"uniqueIndex;not null"`
	ConversationID  string         `json:"conversation_id" gorm:"index;not null"`
	ParentMessageID *string        `json:"parent_message_id" gorm:"index"`
	Role            string         `json:"role" gorm:"size:20;not null"` // user, assistant
	Content         string         `json:"content" gorm:"type:text"`
	MessageMetadata JSONMap        `json:"message_metadata" gorm:"type:json"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// Relationships
	Conversation Conversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationID;references:ConversationID"`
}