package register_routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/from_sso (SSR)
func fromSSO(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Validate the token
	state, msg := validateToken(req.Token, 4)
	if msg != nil {
		return integration.FailedRequest(c, msg, nil)
	}

	// Check if it is really a SSO token
	if !state.SSO {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	return renderUsernameForm(c)
}
