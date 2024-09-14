package register_routes

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/email_code
func checkEmailCode(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate the token
	state, msg := validateToken(req.Token, 3)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// TODO: Add rate limiting
	// Check the email code and stuff
	if state.EmailCode != req.Code {
		return util.FailedRequest(c, localization.ErrorEmailCodeInvalid, nil)
	}

	// Upgrade the token
	if msg := upgradeToken(req.Token, 4); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Render the username creation form
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{}))
}
