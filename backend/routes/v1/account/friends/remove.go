package friends

import (
	"node-backend/database"
	"node-backend/entities/account/properties"
	"node-backend/util"

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
	accId := util.GetAcc(c)
	var friendship properties.Friendship
	if err := database.DBConn.Where("id = ? AND account = ?", req.ID, accId).Take(&friendship).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return util.FailedRequest(c, "not.found", nil)
		}

		return util.FailedRequest(c, "server.error", err)
	}

	// Delete friendship
	if err := database.DBConn.Delete(&friendship).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
