package routes

import (
	"github.com/Liphium/station/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/receive"
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
	if rq.Token != pipes.CurrentNode.Token {
		return integration.InvalidRequest(c, "invalid token")
	}

	receive.HandleMessage("ws", rq.Message)
	return nil
}
