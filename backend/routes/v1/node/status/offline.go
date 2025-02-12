package status

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type offlineRequest struct {
	Token string `json:"token"`
}

func offline(c *fiber.Ctx) error {

	// Parse request
	var req offlineRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get node
	var requested database.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&requested).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Update status
	nodes.TurnOff(&requested, database.StatusStopped)

	if err := database.DBConn.Save(&requested).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
