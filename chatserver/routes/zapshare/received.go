package zapshare_routes

import (
	"fmt"

	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type chunkReceivedRequest struct {
	Id       string `json:"id"`
	Receiver string `json:"receiver"`
	Token    string `json:"token"`
}

func receivedChunk(c *fiber.Ctx) error {

	var req chunkReceivedRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	if req.Id == "" || req.Receiver == "" || req.Token == "" {
		return integration.InvalidRequest(c, "id and token are required")
	}

	transaction, ok := zapshare.GetTransaction(req.Id)
	if !ok || transaction.Token != req.Token {
		return integration.InvalidRequest(c, "invalid id or token")
	}

	if transaction.Token != req.Token {
		return integration.InvalidRequest(c, "invalid token")
	}

	finished, err := transaction.PartReceived(req.Receiver)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprint("invalid receiver: ", err))
	}

	if finished {
		zapshare.CancelTransaction(transaction.Id)

		return c.JSON(fiber.Map{
			"success":  true,
			"complete": true,
		})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"complete": false,
	})
}
