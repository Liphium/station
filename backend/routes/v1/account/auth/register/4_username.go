package register_routes

import (
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/username (SSR)
func checkUsername(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token       string `json:"token"`
		Username    string `json:"username"`
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
	if msg := standards.CheckUsername(req.Username); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Add username and stuff to the state
	state.Mutex.Lock()
	state.Username = req.Username
	state.DisplayName = req.DisplayName
	state.Mutex.Unlock()

	// Upgrade the token
	if msg := upgradeToken(req.Token, 5); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Render the password form
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.RegisterPasswordTitle,
			Style: ssr.TextStyleHeadline,
		},
		ssr.Text{
			Text:  localization.RegisterPasswordRequirements,
			Style: ssr.TextStyleDescription,
		},
		ssr.Input{
			Placeholder: localization.RegisterPasswordPlaceholder,
			Name:        "password",
			Hidden:      true,
		},
		ssr.Input{
			Placeholder: localization.RegisterPasswordConfirmPlaceholder,
			Name:        "confirm_password",
			Hidden:      true,
		},
		ssr.SubmitButton{
			Label: localization.AuthFinishButton,
			Path:  "/account/auth/register/password",
		},
	}))
}