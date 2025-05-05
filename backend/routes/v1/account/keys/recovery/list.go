package recovery_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/recovery/list
func listRecoveryTokens(c *fiber.Ctx) error {

	// Get the account id
	acc, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.FailedRequest(c, localization.ErrorAccountNotFound, err)
	}

	// Get the recovery tokens
	var recoveryTokens []database.RecoveryToken
	if err := database.DBConn.Where("account = ?", acc).Order("created_at").Find(&recoveryTokens).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"list":    recoveryTokens,
	})
}
