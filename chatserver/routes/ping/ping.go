package ping

import (
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

func Pong(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"gateway": integration.Nodes[integration.IdentifierChatNode].NodeId,
		"app":     "chat-node",
	})
}
