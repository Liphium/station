package conversations

type Conversation struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Type uint   `json:"type" gorm:"not null"`
	Data string `json:"data" gorm:"not null"` // Encrypted with the conversation key
}

const TypePrivateMessage = 0
const TypeGroup = 1
