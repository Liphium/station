package pipeshroutes

import (
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
	local := c.Locals("local").(*pipes.LocalNode)
	if event.Token != local.Token {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	local.HandleMessage(pipes.ProtocolWS, event.Message)

	return c.JSON(fiber.Map{
		"success": true,
	})
}
