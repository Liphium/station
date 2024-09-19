package auth_routes

import (
	sso_routes "github.com/Liphium/station/backend/routes/v1/account/auth/sso"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/form (SSR)
func getStartForm(c *fiber.Ctx) error {

	// Redirect to SSO if enabled
	if sso_routes.Enabled {

		// Generate a new SSO token
		token := sso_routes.GenerateSSOToken(c)

		return util.ReturnJSON(c, ssr.RedirectResponse("/account/auth/sso/form", token))
	}

	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Style: ssr.TextStyleHeadline,
			Text:  localization.AuthStartTitle,
		},
		ssr.Input{
			Placeholder: localization.AuthStartEmailPlaceholder,
			Name:        "email",
		},
		ssr.SubmitButton{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/start",
		},
	}))
}
