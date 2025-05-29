package node_action_routes

import (
	"fmt"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type remoteActionRequest struct {
	AppTag string      `json:"app_tag"` // For example: liphium_chat or liphium_spaces
	Sender string      `json:"sender"`  // Domain of the node calling the request
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// Route: /node/actions/send
func sendNodeAction(c *fiber.Ctx) error {

	// Parse the request
	var req remoteActionRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the app by app tag
	var application database.App
	if err := database.DBConn.Where("tag = ?", req.AppTag).Take(&application).Error; err != nil {
		return integration.InvalidRequest(c, "invalid app tag")
	}

	// Get the node with the lowest load to handle the request on
	var lowest database.Node
	if err := database.DBConn.Model(&database.Node{}).Where("app_id = ? AND status = ?", application.ID, database.StatusStarted).Order("load DESC").Take(&lowest).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorNotSetup, nil)
	}

	// Send the remote action to the node
	answer, err := util.PostRequest(util.NodeProtocol+lowest.Domain+"/actions/"+req.Action, fiber.Map{
		"id":     fmt.Sprintf("%d", lowest.ID),
		"token":  lowest.Token,
		"sender": req.Sender,
		"action": req.Action,
		"data":   req.Data,
	})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"answer":  answer,
	})
}
