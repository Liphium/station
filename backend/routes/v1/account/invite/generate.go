package invite_routes

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"

	"github.com/gofiber/fiber/v2"
)

// Route: /account/invite/generate
func generateInvite(c *fiber.Ctx) error {

	// Get invite count of account
	accId := util.GetAcc(c)
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
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Retract one from the invite count of the account
	inviteCount.Count -= 1
	if err := database.DBConn.Save(&inviteCount).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"invite":  invite.ID,
	})
}
