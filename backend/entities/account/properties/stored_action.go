package properties

import "github.com/google/uuid"

type StoredAction struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account   uuid.UUID `json:"-" gorm:"not null"`
	Payload   string    `json:"payload" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
	CreatedAt int64     `json:"-" gorm:"not null,autoCreateTime:milli"`
}

// Authenticated stored actions
type AStoredAction struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account   uuid.UUID `json:"-" gorm:"not null"`
	Payload   string    `json:"payload" gorm:"not null"` // Encrypted payload (encrypted with the account's public key)
	CreatedAt int64     `json:"-" gorm:"not null,autoCreateTime:milli"`
}
