package database

import (
	"regexp"
	"unsafe"

	"github.com/Liphium/station/chatserver/util"
	"github.com/google/uuid"
)

type ConversationToken struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Conversation string `json:"conversation" gorm:"not null,index"` // Conversation ID
	Activated    bool   `json:"activated" gorm:"not null"`          // Whether the token is activated or not
	Token        string `json:"token" gorm:"not null,unique,index"` // Long token required to subscribe to the conversation
	Data         string `json:"data" gorm:"not null"`               // Encrypted data about the user (account id, username, etc.)
	Rank         uint   `json:"rank" gorm:"not null"`
	LastRead     int64  `json:"-" gorm:"not null"` // Last time the user read the conversation
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
const RankModerator = 1 // Can remove/add users
const RankAdmin = 2     // Manages moderators and can delete the conversation

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

func CheckSize(message string) bool {
	return unsafe.Sizeof(message) > 1000*6
}

// Add the extra part added to the conversation id in the messages table (for topics in squares for example)
func WithExtra(conversationId string, extra string) string {
	if extra == "" {
		return conversationId
	}
	return conversationId + "_" + extra
}

// Regex for making sure there are only a-z, A-Z and 1-9 present in the extra part of the conversation id
const extraRegex = "^[a-zA-Z1-9]+$"

// Make sure the extra part isn't weird
func ValidateExtra(extra string) bool {
	if extra == "" {
		return true
	}
	if len(extra) > util.MaxMessageExtra {
		return false
	}
	matched, err := regexp.MatchString(extraRegex, extra)
	return err == nil && matched
}
