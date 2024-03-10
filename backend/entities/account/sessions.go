package account

import (
	"time"
)

type Session struct {
	ID    string `json:"id" gorm:"primaryKey"` //  8 character-long string
	Token string `json:"token" gorm:"unique"`

	Account         string    `json:"account"`
	PermissionLevel uint      `json:"level"`
	Device          string    `json:"device"`
	App             uint      `json:"app"`
	Node            uint      `json:"node"`
	LastUsage       time.Time `json:"last_usage"`
	LastConnection  time.Time `json:"last_connection"` // LastConnection is the last time a new connection was established
}
