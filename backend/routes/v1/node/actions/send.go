package node_action_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/entities/node"
	"github.com/Liphium/station/backend/util"
	"github.com/gofiber/fiber/v2"
)

type remoteActionRequest struct {
	AppTag string      `json:"app_tag"` // For example: liphium_chat or liphium_spaces
	Event  interface{} `json:"event"`   // Event (can be anything)
}

// Route: /node/actions/send
func sendNodeAction(c *fiber.Ctx) error {

	// Parse the request
	var req remoteActionRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the app by app tag
	var application app.App
	if err := database.DBConn.Where("tag = ?", req.AppTag); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the node with the lowest load to handle the request on
	var lowest node.Node
	if err := database.DBConn.Model(&node.Node{}).Where("app_id = ? AND status = ?", application.ID, node.StatusStarted).Order("load DESC").Take(&lowest).Error; err != nil {
		return util.FailedRequest(c, "not.setup", nil)
	}

	// Get public key of node
	res, err := util.PostRequestNoTC(util.NodeProtocol+lowest.Domain+"/pub", fiber.Map{})
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Unpackage the public key
	publicKey, err := util.UnpackageRSAPublicKey(res["pub"].(string))
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Send the remote action to the node
	answer, err := util.PostRequest(publicKey, util.NodeProtocol+lowest.Domain+"/actions/receive", fiber.Map{
		"id":    lowest.ID,
		"token": lowest.Token,
		"event": req.Event,
	})
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"answer":  answer,
	})
}