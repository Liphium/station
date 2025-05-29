package auth_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	login_routes "github.com/Liphium/station/backend/routes/v1/account/auth/login"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/main/ssr"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Route: /account/auth/start
func startAuth(c *fiber.Ctx) error {

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

	// Check if there is an account with this email
	var acc database.Account
	if err := database.DBConn.Where("email = ?", normalizedEmail).Preload("Rank").Take(&acc).Error; err != nil {

		// If the account wasn't found, redirect to register
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(ssr.SuggestResponse(c, localization.ErrorEmailNotFound, ssr.Button{
				Label: localization.AuthStartCreateButton,
				Path:  "/account/auth/register/email",
			}))
		}

		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// If there is an account, redirect to login
	return c.JSON(ssr.RedirectResponse("/account/auth/login/start", login_routes.GenerateLoginToken(acc)))
}
