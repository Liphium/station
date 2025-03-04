package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type listFriendsRequest struct {
	Version int64 `json:"version"`
}

// Route: /account/friends/sync
func listFriends(c *fiber.Ctx) error {

	// Parse to request
	var req listFriendsRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get all of the friends that have been updated
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	var friends []database.Friendship
	if err := database.DBConn.Model(&database.Friendship{}).Where("account = ? AND version > ?", accId, req.Version).Find(&friends).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return new friends
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"friends": friends,
	})
}
