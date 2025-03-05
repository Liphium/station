package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type removeRequest struct {
	ID string `json:"id"`
}

// Route: /account/friends/remove
func removeFriend(c *fiber.Ctx) error {

	// Parse request
	var req removeRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if friendship exists
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return util.InvalidRequest(c)
	}
	var friendship database.Friendship
	if err := database.DBConn.Where("id = ? AND account = ?", req.ID, accId).Take(&friendship).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return util.FailedRequest(c, localization.ErrorFriendNotFound, nil)
		}

		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Get the latest version
	var version int64
	if err := database.DBConn.Model(&database.VaultEntry{}).Select("max(version)").Where("account = ?", accId).Scan(&version).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete the friendship
	friendship.Payload = ""
	friendship.LastPacket = ""
	friendship.Deleted = true
	friendship.Version = version + 1
	if err := database.DBConn.Save(&friendship).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"version": version + 1,
	})
}
