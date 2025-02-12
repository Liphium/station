package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type listFriendsRequest struct {
	After uint64 `json:"after"`
}

// Route: /account/friends/list
func listFriends(c *fiber.Ctx) error {

	var req listFriendsRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get friends list
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	var friends []database.Friendship
	if err := database.DBConn.Model(&database.Friendship{}).Where("account = ? AND updated_at > ?", accId, req.After).Find(&friends).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return friends list
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"friends": friends,
	})
}
