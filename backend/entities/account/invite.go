package account

import "time"

// Invite count for how much individual accounts can generate
type InviteCount struct {
	Account string `gorm:"primaryKey"`
	Count   int    // How many invites can be generated
}

// Invites generated
type Invite struct {
	ID        string    `gorm:"primaryKey"` // Invite token itself
	Creator   string    // Account id of creator
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
