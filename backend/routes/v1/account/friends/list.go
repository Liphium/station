package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

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
	accId := util.GetAcc(c)
	var friends []properties.Friendship
	if err := database.DBConn.Model(&properties.Friendship{}).Where("account = ? AND updated_at > ?", accId, req.After).Find(&friends).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Return friends list
	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"friends": friends,
	})
}
