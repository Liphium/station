package account

import (
	"github.com/Liphium/station/backend/util/auth"
	"github.com/google/uuid"
)

type Authentication struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account uuid.UUID `json:"account"`
	Type    uint      `json:"type"`
	Secret  string    `json:"secret"`
}

const TypePassword = 0
const TypeSSO = 1

// Order to autenticate (0 = first, 1 = second, etc.)
var Order = map[uint]uint{
	TypePassword: 0,
	TypeSSO:      0,
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
	case TypePassword:
		return a.checkPassword(secret, id)
	case TypeSSO:
		return false // TODO: Implement
	}

	return false
}
