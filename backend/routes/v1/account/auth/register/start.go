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
		Email string `json:"email"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	return util.ReturnJSON(c, ssr.RenderResponse(c, ssr.Components{
		ssr.Text{
			Text:  localization.AuthNextStepButton,
			Style: ssr.TextStyleHeadline,
		},
	}))
}
