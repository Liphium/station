package account

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Token string `json:"token" gorm:"unique"`

	Verified        bool      `json:"sync"`
	Account         uuid.UUID `json:"account"`
	PermissionLevel uint      `json:"level"`
	Device          string    `json:"device"`
	App             uint      `json:"app"`
	Node            uint      `json:"node"`
	LastUsage       time.Time `json:"last_usage"`
	LastConnection  time.Time `json:"last_connection"` // LastConnection is the last time a new connection was established
}
