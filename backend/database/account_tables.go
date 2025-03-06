package database

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"`

	Email       string    `json:"email" gorm:"uniqueIndex"`
	Username    string    `json:"username" gorm:"uniqueIndex"`
	DisplayName string    `json:"display_name"`
	RankID      uint      `json:"rank"`
	CreatedAt   time.Time `json:"created_at"`

	Rank           Rank             `json:"-" gorm:"foreignKey:RankID"`
	Authentication []Authentication `json:"-" gorm:"foreignKey:Account"`
	Sessions       []Session        `json:"-" gorm:"foreignKey:Account"`
}

type Session struct {
	ID    uuid.UUID `json:"id" gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"`
	Token string    `json:"token" gorm:"unique;index"`

	Verified        bool      `json:"sync"`
	Account         uuid.UUID `json:"account"`
	PermissionLevel uint      `json:"level"`
	Device          string    `json:"device"`
	App             uint      `json:"app"`
	Node            uint      `json:"node"`
	LastUsage       time.Time `json:"last_usage"`
	LastConnection  time.Time `json:"last_connection"` // LastConnection is the last time a new connection was established
}

type Rank struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Name  string `json:"name"`
	Level uint   `json:"level"`

	Accounts []Account `json:"-" gorm:"foreignKey:RankID"`
}

// * Public keys
type PublicKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`
}

type SignatureKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`
}

// * Symmetric keys
type VaultKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Encrypted with public key and signed
}

type ProfileKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Encrypted with public key and signed
}

// * Keys for safety
type StoredActionKey struct {
	ID  uuid.UUID `json:"id" gorm:"primaryKey"` // Account id
	Key string    `json:"key"`                  // Generated on the server
}

// Invite count for how much individual accounts can generate
type InviteCount struct {
	Account uuid.UUID `gorm:"primaryKey"`
	Count   int       // How many invites can be generated
}

// Invites generated
type Invite struct {
	ID        uuid.UUID `gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"` // Invite token itself
	Creator   uuid.UUID `gorm:"index"`                                           // Account id of creator
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// ! ALL THE DATA IN THIS IS PUBLIC AND CAN BE ACCESED THROUGH /ACCOUNT/PROFILE/GET
type Profile struct {
	ID uuid.UUID `json:"id,omitempty" gorm:"primaryKey"` // Account ID

	//* Picture stuff
	Picture     string `json:"picture,omitempty"`      // File id of the picture
	Container   string `json:"container,omitempty"`    // Attachment container encrypted with profile key
	PictureData string `json:"picture_data,omitempty"` // Profile picture data (zoom, x/y offset) encrypted with profile key

	Data string `json:"data,omitempty"` // Encrypted data (if we need it in the future)
}
