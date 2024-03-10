package rank

import (
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/nodes"

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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check node token
	_, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get rank
	var rank account.Rank
	if database.DBConn.Where("id = ?", req.ID).Find(&rank).Error != nil {
		return util.FailedRequest(c, "server.error", nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"rank":    rank,
	})
}
