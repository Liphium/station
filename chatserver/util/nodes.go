package util

import (
	"errors"
	"log"

	"github.com/Liphium/station/main/integration"
)

type AppToken struct {
	Node   uint   `json:"node"` // Node ID
	Domain string `json:"domain"`
	Token  string `json:"token"`
}

func ConnectToApp(account string, session string, app uint) (AppToken, error) {

	// Get the lowest node
	nodeData := integration.Nodes[integration.IdentifierChatNode]
	log.Println("DATA", account, session)
	res, err := integration.PostRequest("/node/get_lowest", map[string]interface{}{
		"account": account,
		"session": session,
		"app":     app,
		"node":    nodeData.NodeId,
		"token":   nodeData.NodeToken,
	})
	if err != nil {
		return AppToken{}, err
	}

	if !res["success"].(bool) {
		return AppToken{}, errors.New(res["error"].(string))
	}

	return AppToken{
		Node:   uint(res["id"].(float64)),
		Domain: res["domain"].(string),
		Token:  res["token"].(string),
	}, nil
}
