package node

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type thisRequest struct {
	Node  uint   `json:"node"`
	Token string `json:"token"`
}

// Route: /node/this
func this(c *fiber.Ctx) error {

	// Parse request
	var req thisRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get node
	node, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"jwt_secret": util.JWT_SECRET,
		"node":       node.ToEntity(),
	})

}
