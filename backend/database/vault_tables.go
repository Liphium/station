package database

import "github.com/google/uuid"

// All files stored on the server
type CloudFile struct {
	Id        string    `json:"id,omitempty"`   // Format: a-[accountId]-[objectIdentifier]
	Name      string    `json:"name,omitempty"` // File name (encrypted with file key)
	Type      string    `json:"type,omitempty"` // Mime type
	Key       string    `json:"key,omitempty"`  // Encryption key (encrypted with account public key)
	Account   uuid.UUID `json:"account,omitempty"`
	Size      int64     `json:"size,omitempty"`   // In bytes
	Tag       string    `json:"tag,omitempty"`    // Tag for systems such as the library
	System    bool      `json:"system,omitempty"` // If in use by system
	CreatedAt int64     `json:"created,omitempty" gorm:"not null,autoCreateTime:milli"`
}

// All stored actions
type StoredAction struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account   uuid.UUID `json:"-" gorm:"not null"`
	Payload   string    `json:"payload" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
	CreatedAt int64     `json:"-" gorm:"not null,autoCreateTime:milli"`
}

// Authenticated stored actions (stored actions but with a key)
type AStoredAction struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account   uuid.UUID `json:"-" gorm:"not null"`
	Payload   string    `json:"payload" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
	CreatedAt int64     `json:"-" gorm:"not null,autoCreateTime:milli"`
}

// Friend vault
type Friendship struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account    uuid.UUID `json:"account" gorm:"not null"`
	Payload    string    `json:"friend" gorm:"not null"` // Encrypted (with account's public key) friend key + data
	LastPacket string    `json:"-"`                      // When the last packet was received (to prevent replay attacks, encrypted)
	Version    int64     `json:"version" gorm:"index;default:1"`
	Deleted    bool      `json:"deleted" gorm:"index;default:false"`
	UpdatedAt  int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// Vault for all kinds of things (e.g. conversation tokens, etc.)
type VaultEntry struct {
	ID string `json:"id" gorm:"primaryKey"`

	Tag       string    `json:"tag" gorm:"not null;index"` // Tag for the entry (e.g. "conversation")
	Account   uuid.UUID `json:"account" gorm:"not null;index"`
	Payload   string    `json:"payload" gorm:"not null"` // Encrypted (with account's public key) data
	Deleted   bool      `json:"deleted" gorm:"index;default:false"`
	Version   int64     `json:"version" gorm:"default:1;index"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
