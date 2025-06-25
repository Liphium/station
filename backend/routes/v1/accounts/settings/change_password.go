package settings_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type changePasswordRequest struct {
	Current string `json:"current"`
	New     string `json:"new"`
}

// Change the password of an account (Route: /account/settings/change_password)
func changePassword(c *fiber.Ctx) error {

	var req changePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get current password
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}
	var authentication database.Authentication
	if err := database.DBConn.Where("account = ? AND type = ?", accId, database.AuthTypePassword).Take(&authentication).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check password
	if match, err := auth.ComparePasswordAndHash(req.Current, accId, authentication.Secret); err != nil || !match {
		return integration.FailedRequest(c, localization.ErrorPasswordInvalid(8), err)
	}

	// Log out all devices
	// TODO: Disconnect all sessions
	if err := database.DBConn.Where("account = ?", accId).Delete(&database.Session{}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Change password
	hash, err := auth.HashPassword(req.New, accId)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	err = database.DBConn.Model(&database.Authentication{}).Where("account = ? AND type = ?", accId, database.AuthTypePassword).
		Update("secret", hash).Error
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// TODO: Send a mail here in the future (Stuff required: Rate limiting)

	return integration.SuccessfulRequest(c)
}
