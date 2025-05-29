package register_routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/start (from SSR redirect)
func startRegister(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Validate the token
	_, msg := validateToken(req.Token, 1)
	if msg != nil {
		return integration.FailedRequest(c, msg, nil)
	}

	// Upgrade the token to step 2
	if msg := upgradeToken(req.Token, 2); msg != nil {
		return integration.FailedRequest(c, msg, nil)
	}

	// Render the invite input form
	return c.JSON(ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.RegisterInviteTitle,
			Style: ssr.TextStyleHeadline,
		},
		ssr.Input{
			Placeholder: localization.RegisterInvitePlaceholder,
			Name:        "invite",
		},
		ssr.SubmitButton{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/register/invite",
		},
	}))
}
