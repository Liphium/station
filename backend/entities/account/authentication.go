package account

import (
	"strings"

	"github.com/Liphium/station/backend/util/auth"
)

type Authentication struct {
	ID string `json:"id" gorm:"primaryKey"`

	Account string `json:"account"`
	Type    uint   `json:"type"`
	Secret  string `json:"secret"`
}

const TypePassword = 0
const TypeTOTP = 1
const TypeRecoveryCode = 2
const TypePasskey = 3 // Implemented in the future

// Order to autenticate (0 = first, 1 = second, etc.)
var Order = map[uint]uint{
	TypePassword:     0,
	TypePasskey:      5, // Disabled (needs to still be implemented), will eventually be first too
	TypeTOTP:         1,
	TypeRecoveryCode: 1,
}

// Starting step when authenticating
const StartStep = 0

func (a *Authentication) checkPassword(password string, id string) bool {

	match, err := auth.ComparePasswordAndHash(password, id, a.Secret)
	if err != nil {
		return false
	}

	return match
}

func (a *Authentication) Verify(authType uint, secret string, id string) bool {

	if a.Type != authType {
		return false
	}

	switch authType {
	case TypePassword:
		return a.checkPassword(secret, id)
	case TypeTOTP:
		return false // TODO: Implement
	case TypeRecoveryCode:
		return strings.Compare(a.Secret, secret) == 0 // TODO: Implement
	}

	return false
}
