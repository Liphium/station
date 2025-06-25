package register_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/main/integration"
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
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Validate the email
	valid, normalizedEmail := standards.CheckEmail(req.Email)
	if !valid {
		return integration.FailedRequest(c, localization.ErrorEmailInvalid, nil)
	}

	// Make sure there is no other account with this email
	var acc database.Account
	if err := database.DBConn.Where("email = ?", normalizedEmail).Take(&acc).Error; err == nil {
		return integration.FailedRequest(c, localization.ErrorEmailAlreadyInUse, nil)
	}

	// Generate a registration token and redirect to start
	return c.JSON(ssr.RedirectResponse("/account/auth/register/start", GenerateRegisterToken(normalizedEmail)))
}
