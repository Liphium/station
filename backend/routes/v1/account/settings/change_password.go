package settings_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get current password
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	var authentication account.Authentication
	if err := database.DBConn.Where("account = ? AND type = ?", accId, account.TypePassword).Take(&authentication).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check password
	if match, err := auth.ComparePasswordAndHash(req.Current, accId, authentication.Secret); err != nil || !match {
		return util.FailedRequest(c, localization.ErrorPasswordInvalid(8), err)
	}

	// Log out all devices
	// TODO: Disconnect all sessions
	if err := database.DBConn.Where("account = ?", accId).Delete(&account.Session{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Change password
	hash, err := auth.HashPassword(req.New, accId)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	err = database.DBConn.Model(&account.Authentication{}).Where("account = ? AND type = ?", accId, account.TypePassword).
		Update("secret", hash).Error
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// TODO: Send a mail here in the future (Stuff required: Rate limiting)

	return util.SuccessfulRequest(c)
}
