package node

import (
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get node
	node, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return util.InvalidRequest(c)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success":    true,
		"jwt_secret": util.JWT_SECRET,
		"node":       node.ToEntity(),
	})

}
