package requests

import (
	"node-backend/database"
	"node-backend/entities/node"
	"node-backend/util"

	"github.com/gofiber/fiber/v2"
)

type Event struct {
	Sender string                 `json:"sender"` // Sender identifier ("0" for system)
	Name   string                 `json:"name"`
	Data   map[string]interface{} `json:"data"`
}

type channel struct {
	Channel string   `json:"channel"` // Channel name
	Target  []string `json:"target"`  // User IDs to send to (node and user ID for p2p channel)
	Nodes   []string `json:"-"`       // Nodes to send to (only for conversation channel)
}

type message struct {
	Channel channel `json:"channel"`
	Event   Event   `json:"event"`
	NoSelf  bool    `json:"-"` // Whether to send to self (excluded from JSON)
	Local   bool    `json:"-"` // Whether to only send to local clients (excluded from JSON)
}

func SendEventToNode(nodeID uint, account string, event Event) error {

	// Get node
	var receiverNode node.Node
	if err := database.DBConn.Where("id = ?", nodeID).Take(&receiverNode).Error; err != nil {
		return err
	}

	// Get public key of node
	res, err := util.PostRequestNoTC(util.NodeProtocol+receiverNode.Domain+"/pub", fiber.Map{})
	if err != nil {
		return err
	}

	// Unpackage the public key
	publicKey, err := util.UnpackageRSAPublicKey(res["pub"].(string))
	if err != nil {
		return err
	}

	util.PostRequest(publicKey, util.NodeProtocol+receiverNode.Domain+"/adoption/socketless", map[string]interface{}{
		"token": receiverNode.Token,
		"message": message{
			Channel: channel{
				Channel: "p",
				Target:  []string{account},
			},
			Event: event,
		},
	})

	return nil
}
