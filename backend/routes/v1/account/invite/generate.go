package invite_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/invite/generate
func generateInvite(c *fiber.Ctx) error {

	// Get invite count of account
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	var inviteCount account.InviteCount
	if err := database.DBConn.Where("account = ?", accId).Take(&inviteCount).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Check if the account can generate invites
	if inviteCount.Count <= 0 {
		return util.InvalidRequest(c)
	}

	// Generate new invite
	invite := account.Invite{
		ID:      auth.GenerateToken(32),
		Creator: accId,
	}
	if err := database.DBConn.Create(&invite).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Retract one from the invite count of the account
	inviteCount.Count -= 1
	if err := database.DBConn.Save(&inviteCount).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"invite":  invite.ID,
	})
}
