package register_routes

import (
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/username (SSR)
func checkDisplayName(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token       string `json:"token"`
		DisplayName string `json:"display_name"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate the token and stuff
	state, msg := validateToken(req.Token, 4)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Verify display name and username
	if msg := standards.CheckDisplayName(req.DisplayName); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Add username and stuff to the state
	state.Mutex.Lock()
	state.DisplayName = req.DisplayName
	state.Mutex.Unlock()

	// Upgrade the token
	if msg := upgradeToken(req.Token, 5); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Redirect SSO people to step 5 (they don't need a password)
	if state.SSO {
		return util.ReturnJSON(c, ssr.RedirectResponse("/account/auth/register/password", ""))
	}

	// Render the password form
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.RegisterUsernameTitle,
			Style: ssr.TextStyleHeadline,
		},
		ssr.Text{
			Text:  localization.RegisterUsernameDescription,
			Style: ssr.TextStyleDescription,
		},
		ssr.Input{
			Placeholder: localization.RegisterUsernamePlaceholder,
			Name:        "username",
			MaxLength:   standards.MaxDisplayNameLength,
		},
		ssr.SubmitButton{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/register/username",
		},
	}))
}
