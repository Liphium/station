package friends

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account/properties"
	"github.com/Liphium/station/backend/util"
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
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	var friendship properties.Friendship
	if err := database.DBConn.Where("id = ? AND account = ?", req.ID, accId).Take(&friendship).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return util.FailedRequest(c, localization.ErrorFriendNotFound, nil)
		}

		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete friendship
	if err := database.DBConn.Delete(&friendship).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
