package auth_routes

import (
	sso_routes "github.com/Liphium/station/backend/routes/v1/account/auth/sso"
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

		return c.JSON(ssr.RedirectResponse("/account/auth/sso/form", token))
	}

	return c.JSON(ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Style: ssr.TextStyleHeadline,
			Text:  localization.AuthStartTitle,
		},
		ssr.Text{
			Style: ssr.TextStyleDescription,
			Text:  localization.AuthStartDescription,
		},
		ssr.Input{
			Label:       localization.AuthStartEmailLabel,
			Placeholder: localization.AuthStartEmailPlaceholder,
			Name:        "email",
		},
		ssr.SubmitButton{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/start",
		},
	}))
}
