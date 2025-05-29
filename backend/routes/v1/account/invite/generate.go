package invite_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/invite/generate
func generateInvite(c *fiber.Ctx) error {

	// Get invite count of account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Only check for the invite code if the user isn't an admin
	if !verify.InfoLocals(c).HasPermission(verify.PermissionAdmin) {
		var inviteCount database.InviteCount
		if err := database.DBConn.Where("account = ?", accId).Take(&inviteCount).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorInvitesEmpty, err)
		}

		// Check if the account can generate invites
		if inviteCount.Count <= 0 {
			return integration.FailedRequest(c, localization.ErrorInvitesEmpty, nil)
		}

		// Retract one from the invite count of the account
		inviteCount.Count -= 1
		if err := database.DBConn.Save(&inviteCount).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	// Generate new invite
	invite := database.Invite{
		Creator: accId,
	}
	if err := database.DBConn.Create(&invite).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"invite":  invite.ID.String(),
	})
}
