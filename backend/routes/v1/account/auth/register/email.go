package register_routes

import (
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/auth/register/email
func checkEmail(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Email string `json:"email"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Validate the email
	valid, normalizedEmail := standards.CheckEmail(req.Email)
	if !valid {
		return util.FailedRequest(c, localization.ErrorEmailInvalid, nil)
	}

	// Generate a registration token and redirect to start
	return util.ReturnJSON(c, ssr.RedirectResponse("/account/auth/register/start", GenerateRegisterToken(normalizedEmail)))
}
