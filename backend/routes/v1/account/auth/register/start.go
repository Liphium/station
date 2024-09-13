package register_routes

import (
	"github.com/Liphium/station/backend/util"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate the token
	state, msg := validateToken(req.Token, 1)
	if msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Render the invite input form
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.RegisterInviteTitle,
			Style: ssr.TextStyleHeadline,
		},
		ssr.Input{
			Placeholder: localization.AuthStartEmailPlaceholder,
			Value:       state.Email,
			Name:        "email",
		},
		ssr.Input{
			Placeholder: localization.RegisterInvitePlaceholder,
			Name:        "invite",
		},
		ssr.Button{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/register/invite",
		},
		ssr.Button{
			Label: localization.AuthResendEmailButton,
			Path:  "/account/auth/register/resend_email",
		},
	}))
}
