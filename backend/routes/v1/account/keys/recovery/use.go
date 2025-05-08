package recovery_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/recovery/use
func useRecoveryToken(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		SessionToken  string `json:"session_token"`
		RecoveryToken string `json:"recovery_token"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the session
	var session database.Session
	if err := database.DBConn.Where("token = ?", req.SessionToken).Take(&session).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Check if the recovery token is valid
	if err := database.DBConn.Where("account = ? AND token = ?", session.Account, req.RecoveryToken).Take(&database.RecoveryToken{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRecoveryToken, err)
	}

	// Use the token to verify the session
	session.Verified = true
	if err := database.DBConn.Save(&session).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete the recovery token
	if err := database.DBConn.Where("account = ? AND token = ?", session.Account, req.RecoveryToken).Delete(&database.RecoveryToken{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete any key requests left in the session
	database.DBConn.Where("session = ?", session.ID).Delete(&database.KeyRequest{})

	// Return the payload to the client
	return util.SuccessfulRequest(c)
}
