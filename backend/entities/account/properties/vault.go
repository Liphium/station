package properties

import "github.com/google/uuid"

// Friend vault
type Friendship struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account    uuid.UUID `json:"account" gorm:"not null"`
	Hash       string    `json:"hash" gorm:"not null"`
	Payload    string    `json:"friend" gorm:"not null"` // Encrypted (with account's public key) friend key + data
	LastSent   string    `json:"-"`                      // When the last packet was sent (to prevent replay attacks, encrypted)
	LastPacket string    `json:"-"`                      // When the last packet was received (to prevent replay attacks, encrypted)
	UpdatedAt  int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// Vault for all kinds of things (e.g. conversation tokens, etc.)
type VaultEntry struct {
	ID string `json:"id" gorm:"primaryKey"`

	Tag       string    `json:"tag" gorm:"not null"` // Tag for the entry (e.g. "conversation")
	Account   uuid.UUID `json:"account" gorm:"not null"`
	Payload   string    `json:"payload" gorm:"not null"` // Encrypted (with account's public key) data
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
