package pipeshroutes

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func adoptionRouter(router fiber.Router, local *pipes.LocalNode) {
	router.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {

			// Check if the request has a token
			token := c.Get("Sec-WebSocket-Protocol")

			// Adopt node
			node, err := local.ReceiveWSAdoption(token)
			if err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			// Set the token as a local variable
			c.Locals("ws", true)
			c.Locals("node", node)
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	router.Get("/", websocket.New(func(c *websocket.Conn) {
		adoptionWs(c, local)
	}))
}

func adoptionWs(conn *websocket.Conn, local *pipes.LocalNode) {
	node := conn.Locals("node").(pipes.Node)

	defer func() {

		// Disconnect node
		local.RemoveNodeWS(node.ID)
		pipeshandler.CurrentConfig.NodeDisconnectHandler(node)
		conn.Close()
	}()

	for {
		// Read message as text
		mtype, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if mtype == websocket.TextMessage {

			// Pass message to pipes
			local.ReceiveWS(msg)
		}
	}

}
