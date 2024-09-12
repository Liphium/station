package auth_routes

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/form (SSR)
func getStartForm(c *fiber.Ctx) error {
	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Style: ssr.TextStyleHeadline,
			Text:  localization.AuthStartTitle,
		},
		ssr.Input{
			Placeholder: localization.AuthStartEmailPlaceholder,
			Name:        "email",
		},
		ssr.Button{
			Label: localization.AuthNextStepButton,
			Path:  "/account/auth/start",
		},
	}))
}
