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
	Hash       string    `json:"hash" gorm:"not null"`
	Payload    string    `json:"friend" gorm:"not null"` // Encrypted (with account's public key) friend key + data
	LastPacket string    `json:"-"`                      // When the last packet was received (to prevent replay attacks, encrypted)
	UpdatedAt  int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// All deletions that happen in the vault
type VaultDeletion struct {
	ID      uuid.UUID `gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"`
	Account uuid.UUID `gorm:"type:uuid;index"`

	// Data about the deletion
	Tag     string `gorm:"index"`
	Version int64  `gorm:"not null;default:0;index"` // Version of the tag
	Entry   string `gorm:"not null"`                 // Id of the deleted entry
}

// Vault for all kinds of things (e.g. conversation tokens, etc.)
type VaultEntry struct {
	ID string `json:"id" gorm:"primaryKey"`

	Tag       string    `json:"tag" gorm:"not null;index"` // Tag for the entry (e.g. "conversation")
	Account   uuid.UUID `json:"account" gorm:"not null;index"`
	Payload   string    `json:"payload" gorm:"not null"` // Encrypted (with account's public key) data
	Version   int64     `json:"-" gorm:"default:0;index"`
	UpdatedAt int64     `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
