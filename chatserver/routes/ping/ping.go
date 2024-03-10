package ping

import (
	"github.com/Liphium/station/chatserver/util"
	"github.com/gofiber/fiber/v2"
)

func Pong(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"gateway": util.NODE_ID,
		"app":     "chat-node",
	})
}
