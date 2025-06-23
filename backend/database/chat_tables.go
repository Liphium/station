package database

import (
	"github.com/google/uuid"
)

type ConversationToken struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Conversation string `json:"conversation" gorm:"not null,index"` // Conversation ID
	Activated    bool   `json:"activated" gorm:"not null"`          // Whether the token is activated or not
	Token        string `json:"token" gorm:"not null,unique,index"` // Long token required to subscribe to the conversation
	Data         string `json:"data" gorm:"not null"`               // Encrypted data about the user (account id, username, etc.)
	Rank         uint   `json:"rank" gorm:"not null"`
	Reads        string `json:"reads"`

	// For synchronization data (unrelated to database model)
	LastSync int64 `json:"-" gorm:"-"`
}

func (t *ConversationToken) ToSent() SentConversationToken {
	return SentConversationToken{
		ID:           t.ID,
		Token:        t.Token,
		Conversation: t.Conversation,
	}
}

type SentConversationToken struct {
	ID           string `json:"id"`
	Token        string `json:"token"`
	Conversation string `json:"conv"`
	LastMessage  int64  `json:"time,omitempty"`
}

// * Ranks
const RankUser = 0
const RankModerator = 100 // Can remove/add users
const RankAdmin = 200     // Manages moderators and can delete the conversation

type Conversation struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Type    uint   `json:"type" gorm:"not null"`
	Version int64  `json:"updated" gorm:"not null,default:1"` // The version of the conversation (used to track updates to it)
	Data    string `json:"data" gorm:"not null"`              // Encrypted with the conversation key
}

const ConvTypePrivateMessage = 0
const ConvTypeGroup = 1
const ConvTypeSquare = 2 // :)

type Message struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"`

	Conversation string `json:"cv" gorm:"not null,index"`
	Creation     int64  `json:"ct" gorm:"index"`    // Unix timestamp
	Data         string `json:"dt" gorm:"not null"` // Encrypted data
	Edited       bool   `json:"ed" gorm:"not null"` // Edited flag
	Sender       string `json:"sr" gorm:"not null"` // Sender ID (of conversation token)
}
