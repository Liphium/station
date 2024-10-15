package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/me
func me(c *fiber.Ctx) error {

	// Get session
	sessionId, err := verify.InfoLocals(c).GetSessionUUID()
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	var session database.Session
	if database.DBConn.Where(&database.Session{ID: sessionId}).Take(&session).Error != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	var acc database.Account
	if err := database.DBConn.Where(&database.Account{ID: session.Account}).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get all valid permissions the account has
	perms := []string{}
	info := verify.InfoLocals(c)
	for name := range verify.Permissions {
		if info.HasPermission(name) {
			perms = append(perms, name)
		}
	}

	// Get all the ranks
	var ranks []database.Rank
	if err := database.DBConn.Find(&ranks).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Retrun details
	return util.ReturnJSON(c, fiber.Map{
		"success":     true,
		"account":     acc,
		"permissions": perms,
		"ranks":       ranks,
	})
}
