package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/me
func me(c *fiber.Ctx) error {

	// Get session
	sessionId, err := util.GetSession(c)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	var session account.Session
	if database.DBConn.Where(&account.Session{ID: sessionId}).Take(&session).Error != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where(&account.Account{ID: session.Account}).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get all valid permissions the account has
	perms := []string{}
	for name := range util.Permissions {
		if util.Permission(c, name) {
			perms = append(perms, name)
		}
	}

	// Retrun details
	return util.ReturnJSON(c, fiber.Map{
		"success":     true,
		"account":     acc,
		"permissions": perms,
	})
}
