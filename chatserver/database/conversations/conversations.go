package conversations

type Conversation struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Type    uint   `json:"type" gorm:"not null"`
	Version int64  `json:"updated" gorm:"not null,default:1"` // The version of the conversation (used to track updates to it)
	Data    string `json:"data" gorm:"not null"`              // Encrypted with the conversation key
}

const TypePrivateMessage = 0
const TypeGroup = 1
