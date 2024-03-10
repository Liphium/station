package node

import (
	"log"
	"node-backend/util"
	"node-backend/util/requests"

	"github.com/gofiber/fiber/v2"
)

type sendRequest struct {
	Node    uint   `json:"node"`
	Account string `json:"account"`
	Message string `json:"message"`
}

func sendToNode(c *fiber.Ctx) error {

	// Parse request
	var req sendRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Send event
	if err := requests.SendEventToNode(req.Node, req.Account, requests.Event{
		Sender: "0",
		Name:   "message",
		Data: map[string]interface{}{
			"message": req.Message,
		},
	}); err != nil {
		return util.FailedRequest(c, "node.error", err)
	}

	log.Println("sent to pipes")

	return nil
}
