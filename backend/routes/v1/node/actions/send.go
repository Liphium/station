package node_action_routes

import (
	"fmt"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
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
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Get the app by app tag
	var application database.App
	if err := database.DBConn.Where("tag = ?", req.AppTag).Take(&application).Error; err != nil {
		return util.InvalidRequest(c)
	}

	// Get the node with the lowest load to handle the request on
	var lowest database.Node
	if err := database.DBConn.Model(&database.Node{}).Where("app_id = ? AND status = ?", application.ID, database.StatusStarted).Order("load DESC").Take(&lowest).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorNotSetup, nil)
	}

	// Get public key of node
	res, err := util.PostRequestNoTC(util.NodeProtocol+lowest.Domain+"/pub", fiber.Map{})
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Unpackage the public key
	publicKey, err := util.UnpackageRSAPublicKey(res["pub"].(string))
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send the remote action to the node
	answer, err := util.PostRequest(publicKey, util.NodeProtocol+lowest.Domain+"/actions/"+req.Action, fiber.Map{
		"id":     fmt.Sprintf("%d", lowest.ID),
		"token":  lowest.Token,
		"sender": req.Sender,
		"action": req.Action,
		"data":   req.Data,
	})
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"answer":  answer,
	})
}
