package status

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type offlineRequest struct {
	Token string `json:"token"`
}

func offline(c *fiber.Ctx) error {

	// Parse request
	var req offlineRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get node
	var requested database.Node
	if err := database.DBConn.Where("token = ?", req.Token).Take(&requested).Error; err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Update status
	nodes.TurnOff(&requested, database.StatusStopped)

	if err := database.DBConn.Save(&requested).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
