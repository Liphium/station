package account

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type getRequestNode struct {
	ID    string `json:"id"`
	Node  uint   `json:"node"`
	Token string `json:"token"`
}

// Route: /account/get_node
func getAccountNode(c *fiber.Ctx) error {

	// Parse request
	var req getRequestNode
	if err := util.BodyParser(c, &req); err != nil {
		util.Log.Println(err)
		return util.InvalidRequest(c)
	}

	_, err := nodes.Node(req.Node, req.Token)
	if err != nil {
		return util.InvalidRequest(c)
	}

	// Get account
	var acc database.Account
	if err := database.DBConn.Select("username", "tag").Where("id = ?", req.ID).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	var pub database.PublicKey
	if err := database.DBConn.Select("key").Where("id = ?", req.ID).Take(&pub).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	var signaturePub database.SignatureKey
	if err := database.DBConn.Select("key").Where("id = ?", req.ID).Take(&signaturePub).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"name":    acc.Username,
		"sg":      signaturePub.Key,
		"pub":     pub.Key,
	})
}
