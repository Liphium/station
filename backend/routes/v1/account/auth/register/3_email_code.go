package register_routes

import (
	"time"

	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/email_code (SSR)
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

	// Rate limit the entering of codes
	if !ratelimitHandler(state, 3, time.Second*5) {
		return util.FailedRequest(c, localization.ErrorAuthRatelimit, nil)
	}

	// Check the email code and stuff
	if state.EmailCode != req.Code {
		return util.FailedRequest(c, localization.ErrorEmailCodeInvalid, nil)
	}

	// Upgrade the token
	if msg := upgradeToken(req.Token, 4); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	return renderUsernameForm(c)
}

// Render the username creation form
func renderUsernameForm(c *fiber.Ctx) error {
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.RegisterDisplayNameTitle,
			Style: ssr.TextStyleHeadline,
		},
		ssr.Text{
			Text:  localization.RegisterDisplayNameDescription,
			Style: ssr.TextStyleDescription,
		},
		ssr.Input{
			Placeholder: localization.RegisterDisplayNamePlaceholder,
			Name:        "display_name",
			MaxLength:   standards.MaxUsernameLength,
		},
		ssr.SubmitButton{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/register/display_name",
		},
	}))
}
