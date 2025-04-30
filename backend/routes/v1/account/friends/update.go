package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type updateRequest struct {
	Entry   string `json:"entry"`
	Payload string `json:"payload"`
}

// Route: /account/friend/update
func updateFriend(c *fiber.Ctx) error {

	// Parse the request
	var req updateRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the current account id
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get the friendship
	var entry database.Friendship
	if err := database.DBConn.Where("id = ? AND account = ?", req.Entry, accId).Take(&entry).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the latest version in the friendship
	var version int64
	if err := database.DBConn.Model(&database.Friendship{}).Select("max(version)").Where("account = ?", accId).Scan(&version).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Update the entry to the newest version
	entry.Payload = req.Payload
	entry.Version = version + 1
	if err := database.DBConn.Save(&entry).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"version": version + 1,
	})
}
