package routes

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/gofiber/fiber/v2"
)

type socketlessRq struct {
	This    string        `json:"this"`
	Token   string        `json:"token"`
	Message pipes.Message `json:"message"`
}

func socketlessEvent(c *fiber.Ctx) error {

	// Parse request
	var rq socketlessRq
	if err := c.BodyParser(&rq); err != nil {
		return integration.InvalidRequest(c, "invalid request body, err: "+err.Error())
	}

	// Check token
	if rq.Token != caching.Node.Token {
		return integration.InvalidRequest(c, "invalid token")
	}

	caching.Node.HandleMessage("ws", rq.Message)
	return nil
}
