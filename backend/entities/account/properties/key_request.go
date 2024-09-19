package properties

import "github.com/google/uuid"

type KeyRequest struct {
	Session   uuid.UUID `gorm:"primaryKey" json:"session"`
	Account   uuid.UUID `gorm:"not null" json:"-"`
	Key       string    `json:"pub"`       // Public key of the session requesting it
	Signature string    `json:"signature"` // Signature of the session requesting it
	Payload   string    `json:"payload"`   // Encrypted payload (from the session sending it)
	CreatedAt int64     `json:"creation" gorm:"not null,autoCreateTime:milli"`
}
