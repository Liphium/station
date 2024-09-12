package login_routes

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/login/start (redirect to using SSR)
func startLogin(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check the token
	if msg := testTokenAndRatelimit(req.Token, 1); msg != nil {
		return util.FailedRequest(c, msg, nil)
	}

	// Render the password login page
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Style: ssr.TextStyleHeadline,
			Text:  localization.LoginPasswordTitle,
		},
		ssr.Input{
			Placeholder: localization.LoginPasswordPlaceholder,
			Hidden:      true,
			Name:        "password",
		},
		ssr.Button{
			Label: localization.AuthSubmitButton,
			Path:  "/account/auth/login/password",
		},
	}))
}
