package conversations

type ConversationToken struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Conversation string `json:"conversation" gorm:"not null"` // Conversation ID
	Activated    bool   `json:"activated" gorm:"not null"`    // Whether the token is activated or not
	Token        string `json:"token" gorm:"not null,unique"` // Long token required to subscribe to the conversation
	Data         string `json:"data" gorm:"not null"`         // Encrypted data about the user (account id, username, etc.)
	Rank         uint   `json:"rank" gorm:"not null"`
	LastRead     int64  `json:"-" gorm:"not null"`    // Last time the user read the conversation
	Node         int64  `json:"node" gorm:"not null"` // Node ID
}

type SentConversationToken struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// * Permissions
const MinRankManageMembers = RankModerator
const MinRankChangeConversationDetails = RankModerator
const MinRankManageModerators = RankAdmin
const MinRankDeleteConversation = RankAdmin

// * Ranks
const RankUser = 0
const RankModerator = 1 // Can remove/add users
const RankAdmin = 2     // Manages moderators and can delete the conversation
