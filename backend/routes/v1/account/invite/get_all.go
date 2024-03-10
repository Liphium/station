package invite_routes

import (
	"errors"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Route: /account/invite/get_all
func getAllInformation(c *fiber.Ctx) error {

	// Retrieve all the information
	accId := util.GetAcc(c)

	var invitesGenerated []string
	if err := database.DBConn.Model(&account.Invite{}).Where("creator = ?", accId).Limit(30).Order("created_at DESC").Select("id").Scan(&invitesGenerated).Error; err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	var inviteCount account.InviteCount
	err := database.DBConn.Where("account = ?", accId).Take(&inviteCount).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		inviteCount.Count = 0
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"invites": invitesGenerated,
		"count":   inviteCount.Count,
	})
}
