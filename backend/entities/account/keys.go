package account

import "github.com/google/uuid"

type PublicKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`
}

type ProfileKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Encrypted with private key
}

type StoredActionKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Generated on the server
}

type SignatureKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Encrypted with private key
}
