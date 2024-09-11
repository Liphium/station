package auth_routes

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type startRequest struct {
	Email string `json:"email"`
}

// Route: /account/auth/start
func startAuth(c *fiber.Ctx) error {

	// Parse the request
	var req startRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if there is an account with this email
	var acc account.Account
	if err := database.DBConn.Where("email = ?", req.Email).Take(&acc).Error; err != nil {

		// If the account wasn't found, redirect to register
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.ReturnJSON(c, fiber.Map{
				"success":  true,
				"redirect": "/account/auth/register",
			})
		}

		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// If there is an account, redirect to login
	return util.SuccessfulRequest(c)
}
