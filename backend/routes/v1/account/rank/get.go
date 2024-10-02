package rank

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check node token
	_, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get rank
	var rank database.Rank
	if database.DBConn.Where("id = ?", req.ID).Find(&rank).Error != nil {
		return util.FailedRequest(c, localization.ErrorServer, nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"rank":    rank,
	})
}
