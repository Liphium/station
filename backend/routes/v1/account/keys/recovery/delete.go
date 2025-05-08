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
	accId := verify.InfoLocals(c).GetAccount()

	// Get the token
	if err := database.DBConn.Where("account = ? AND token = ?", accId, req.Token).Take(&database.RecoveryToken{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRecoveryToken, err)
	}

	// Delete the token
	if err := database.DBConn.Where("account = ? AND token = ?", accId, req.Token).Delete(&database.RecoveryToken{}).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
