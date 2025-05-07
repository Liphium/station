package recovery_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/recovery/clear
func deleteRecoveryToken(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Token string `json:"token"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequestContent, err)
	}
	acc, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.FailedRequest(c, localization.ErrorAccountNotFound, err)
	}

	// Get the token
	if err := database.DBConn.Where(&database.RecoveryToken{Account: acc, Token: req.Token}).Take(&database.RecoveryToken{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRecoveryToken, err)
	}

	// Delete the token
	if err := database.DBConn.Where(&database.RecoveryToken{Account: acc, Token: req.Token}).Delete(&database.RecoveryToken{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
