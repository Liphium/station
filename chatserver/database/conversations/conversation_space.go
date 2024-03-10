package conversations

type ConversationSpace struct {
	ID   uint   `json:"id" gorm:"primaryKey"` // Conversation ID
	Node uint   `json:"node"`                 // Space node (where the Space is running)
	Data string `json:"data"`                 // Space data
}

type SpaceData struct {
	Clients []string `json:"clients"` // List of clients in the space
	Info    string   `json:"info"`    // Space info (encrypted using Space key)
}
