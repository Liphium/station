package conversations

import (
	"unsafe"

	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"`

	Conversation string `json:"cv" gorm:"not null"`
	Creation     int64  `json:"ct"`                 // Unix timestamp (SET BY THE CLIENT, EXTREMELY IMPORTANT FOR SIGNATURES)
	Data         string `json:"dt" gorm:"not null"` // Encrypted data
	Edited       bool   `json:"ed" gorm:"not null"` // Edited flag
	Sender       string `json:"sr" gorm:"not null"` // Sender ID (of conversation token)
}

func CheckSize(message string) bool {
	return unsafe.Sizeof(message) > 1000*6
}
