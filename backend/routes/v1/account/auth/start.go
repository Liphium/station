package auth_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	login_routes "github.com/Liphium/station/backend/routes/v1/account/auth/login"
	"github.com/Liphium/station/backend/util"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if there is an account with this email
	var acc account.Account
	if err := database.DBConn.Where("email = ?", req.Email).Preload("Rank").Take(&acc).Error; err != nil {

		// If the account wasn't found, redirect to register
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.ReturnJSON(c, ssr.SuggestResponse(c, localization.ErrorEmailNotFound, ssr.Button{
				Label: localization.AuthStartCreateButton,
				Path:  "/account/auth/register/start",
			}))
		}

		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// If there is an account, redirect to login
	return util.ReturnJSON(c, ssr.RedirectResponse("/account/auth/login/start", login_routes.GenerateLoginToken(acc)))
}