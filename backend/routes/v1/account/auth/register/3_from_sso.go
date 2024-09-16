package register_routes

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/from_sso (SSR)
func fromSSO(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate the token
	state, msg := validateToken(req.Token, 4)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Check if it is really a SSO token
	if !state.SSO {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	return renderUsernameForm(c)
}
