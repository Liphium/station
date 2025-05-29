package rank

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type getRequest struct {

	// Rank ID
	ID uint `json:"id"`

	// Node data
	Node  uint   `json:"node"`  // Node ID
	Token string `json:"token"` // Node token
}

func getRank(c *fiber.Ctx) error {

	// Parse request
	var req getRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Check node token
	_, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "invalid node")
	}

	// Get rank
	var rank database.Rank
	if database.DBConn.Where("id = ?", req.ID).Find(&rank).Error != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"rank":    rank,
	})
}
