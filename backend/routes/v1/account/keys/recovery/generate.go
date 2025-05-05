package recovery_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/keys/recovery/generate
func generateRecoveryToken(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Data string `json:"data"` // Encrypted data
	}
	if err := util.BodyParser(c, &req); err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequestContent, err)
	}
	acc, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.FailedRequest(c, localization.ErrorAccountNotFound, err)
	}

	// Make sure there aren't too many recovery tokens
	var count int64
	if err := database.DBConn.Model(&database.RecoveryToken{}).Where("account = ?", acc).Count(&count).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}
	if count >= MaxRecoveryTokens {
		return util.FailedRequest(c, localization.ErrorRecoveryTokenLimitReached(MaxRecoveryTokens), nil)
	}

	// Generate a new recovery token for the account
	token := database.RecoveryToken{
		Account: acc,
		Token:   auth.GenerateToken(40),
		Data:    req.Data,
	}
	if err := database.DBConn.Create(&token).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return the generated token
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   token.Token,
	})
}
