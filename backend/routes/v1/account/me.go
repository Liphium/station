package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Route: /account/me
func me(c *fiber.Ctx) error {

	// Get session
	sessionId := util.GetSession(c)

	var session account.Session
	if database.DBConn.Where(&account.Session{ID: sessionId}).Take(&session).Error != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	var acc account.Account
	if err := database.DBConn.Where(&account.Account{ID: session.Account}).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Get account
	var profile properties.Profile = properties.Profile{}
	if err := database.DBConn.Where("id = ?", acc.ID).Take(&profile).Error; err != nil && err != gorm.ErrRecordNotFound {
		return util.FailedRequest(c, "server.error", err)
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
		"success":      true,
		"account":      acc,
		"permissions":  perms,
		"display_name": profile.DisplayName,
	})
}
