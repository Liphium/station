package recovery_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/recovery/get
func getRecoveryData(c *fiber.Ctx) error {

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
	var recoveryToken database.RecoveryToken
	if err := database.DBConn.Where("account = ? AND token = ?", session.Account, req.RecoveryToken).Take(&recoveryToken).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRecoveryToken, err)
	}

	// Return the payload to the client
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"data":    recoveryToken.Data,
	})
}
