package database

import (
	"github.com/Liphium/station/backend/util/auth"
	"github.com/google/uuid"
)

type Authentication struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey,type:uuid;default:uuid_generate_v4()"`

	Account uuid.UUID `json:"account"`
	Type    uint      `json:"type"`
	Secret  string    `json:"secret"`
}

const AuthTypePassword = 0
const AuthTypeSSO = 1

// Order to autenticate (0 = first, 1 = second, etc.)
var Order = map[uint]uint{
	AuthTypePassword: 0,
	AuthTypeSSO:      0,
}

// Starting step when authenticating
const StartStep = 0

func (a *Authentication) checkPassword(password string, id uuid.UUID) bool {

	match, err := auth.ComparePasswordAndHash(password, id, a.Secret)
	if err != nil {
		return false
	}

	return match
}

func (a *Authentication) Verify(authType uint, secret string, id uuid.UUID) bool {

	if a.Type != authType {
		return false
	}

	switch authType {
	case AuthTypePassword:
		return a.checkPassword(secret, id)
	case AuthTypeSSO:
		return false // TODO: Implement
	}

	return false
}

type KeyRequest struct {
	Session   uuid.UUID `gorm:"primaryKey" json:"session"`
	Account   uuid.UUID `gorm:"not null" json:"-"`
	Key       string    `json:"pub"`       // Public key of the session requesting it
	Signature string    `json:"signature"` // Signature of the session requesting it
	Payload   string    `json:"payload"`   // Encrypted payload (from the session sending it)
	CreatedAt int64     `json:"creation" gorm:"not null,autoCreateTime:milli"`
}

// Recovery tokens for the encryption keys
type RecoveryToken struct {
	Account   uuid.UUID `json:"-" gorm:"index"`
	Token     string    `json:"-"`
	Data      string    `json:"-"` // All of the keys (encrypted with the token)
	CreatedAt int64     `json:"creation" gorm:"autoCreateTime:milli"`
}
