package routes

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/gofiber/fiber/v2"
)

type socketlessEvent struct {
	Token   string        `json:"token"`
	Message pipes.Message `json:"message"`
}

func socketless(c *fiber.Ctx) error {

	// Parse request
	var event socketlessEvent
	if err := c.BodyParser(&event); err != nil {
		return err
	}

	// Check token
	if event.Token != caching.CSNode.Token {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	caching.CSNode.HandleMessage(pipes.ProtocolWS, event.Message)

	return integration.SuccessfulRequest(c)
}
