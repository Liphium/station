package account

import (
	"time"

	"github.com/google/uuid"
)

// Invite count for how much individual accounts can generate
type InviteCount struct {
	Account uuid.UUID `gorm:"primaryKey"`
	Count   int       // How many invites can be generated
}

// Invites generated
type Invite struct {
	ID        uuid.UUID `gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"` // Invite token itself
	Creator   uuid.UUID // Account id of creator
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
