package account

import "github.com/google/uuid"

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
