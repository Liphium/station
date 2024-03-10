package connect

import (
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/nodes"

	"github.com/gofiber/fiber/v2"
)

type disconnectRequest struct {
	Node      uint   `json:"node"`
	NodeToken string `json:"token"`
	Session   string `json:"session"`
}

// Route: /node/disconnect
func Disconnect(c *fiber.Ctx) error {

	// Parse request
	var req disconnectRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check node
	_, err := nodes.Node(req.Node, req.NodeToken)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Disconnect account
	if database.DBConn.Model(&account.Session{}).Where("id = ?", req.Session).Update("node", 0).Error != nil {
		log.Println("Failed to disconnect account", req.Session)
		return util.FailedRequest(c, "server.error", nil)
	}

	return util.SuccessfulRequest(c)
}
