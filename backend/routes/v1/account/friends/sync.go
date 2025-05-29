package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type listFriendsRequest struct {
	Version int64 `json:"version"`
}

// Route: /account/friends/sync
func syncFriends(c *fiber.Ctx) error {

	// Parse to request
	var req listFriendsRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get all of the friends that have been updated
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}
	var friends []database.Friendship
	if err := database.DBConn.Model(&database.Friendship{}).Where("account = ? AND version > ?", accId, req.Version).Find(&friends).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return new friends
	return c.JSON(fiber.Map{
		"success": true,
		"friends": friends,
	})
}
