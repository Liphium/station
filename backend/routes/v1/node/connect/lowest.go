package connect

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/node"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LowestUsageRequest struct {
	Account string `json:"account"`
	Session string `json:"session"`
	Cluster uint   `json:"cluster"`
	App     uint   `json:"app"`
	Node    uint   `json:"node"`  // Node ID
	Token   string `json:"token"` // Node token
}

// Route: /node/get_lowest
func GetLowest(c *fiber.Ctx) error {

	// Parse request
	var req LowestUsageRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Parse account id from request
	id, err := uuid.Parse(req.Account)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Check node
	_, err = nodes.Node(req.Node, req.Token)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get lowest load node
	var lowest node.Node
	search := node.Node{
		ClusterID: req.Cluster,
		AppID:     req.App,
		Status:    node.StatusStarted,
	}

	if err := database.DBConn.Model(&node.Node{}).Where(&search).Order("load DESC").Take(&lowest).Error; err != nil {
		return util.FailedRequest(c, "not.setup", nil)
	}

	// Ping node (to see if it's online)
	if err := lowest.SendPing(); err != nil {

		// Set the node to error
		nodes.TurnOff(&lowest, node.StatusError)
		return util.FailedRequest(c, util.ErrorNode, err)
	}

	// Generate a jwt token for the node
	token, err := util.ConnectionToken(id, req.Session, lowest.ID)
	if err != nil {
		return util.FailedRequest(c, util.ErrorServer, err)
	}

	// Save node
	if err := database.DBConn.Save(&lowest).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"domain":  lowest.Domain,
		"id":      lowest.ID,
		"token":   token,
	})
}
