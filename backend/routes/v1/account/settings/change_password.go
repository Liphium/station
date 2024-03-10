package settings_routes

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"

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
	accId := util.GetAcc(c)
	var authentication account.Authentication
	if err := database.DBConn.Where("account = ? AND type = ?", accId, account.TypePassword).Take(&authentication).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Check password
	if match, err := auth.ComparePasswordAndHash(req.Current, accId, authentication.Secret); err != nil || !match {
		return util.FailedRequest(c, util.PasswordInvalid, err)
	}

	// Log out all devices
	// TODO: Disconnect all sessions
	if err := database.DBConn.Where("account = ?", accId).Delete(&account.Session{}).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Change password
	hash, err := auth.HashPassword(req.New, accId)
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	err = database.DBConn.Model(&account.Authentication{}).Where("account = ? AND type = ?", accId, account.TypePassword).
		Update("secret", hash).Error
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// TODO: Send a mail here in the future (Stuff required: Rate limiting)

	return util.SuccessfulRequest(c)
}
