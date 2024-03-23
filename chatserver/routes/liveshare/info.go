package liveshare_routes

import (
	"github.com/Liphium/station/chatserver/liveshare"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type infoRequest struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

func transactionInfo(c *fiber.Ctx) error {

	var req infoRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	transaction, ok := liveshare.GetTransaction(req.Id)
	if !ok {
		return integration.InvalidRequest(c, "transaction not found")
	}

	if transaction.Token != req.Token {
		return integration.InvalidRequest(c, "invalid token")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"name":    transaction.FileName,
		"size":    transaction.FileSize,
	})
}
