package account

import "github.com/google/uuid"

//* Public keys
type PublicKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`
}

type SignatureKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`
}

//* Symmetric keys
type VaultKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Encrypted with public key and signed
}

type ProfileKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Encrypted with public key and signed
}

//* Keys for safety
type StoredActionKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Generated on the server
}
