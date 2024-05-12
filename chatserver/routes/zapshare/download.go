package zapshare_routes

import (
	"strconv"

	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

func downloadChunk(c *fiber.Ctx) error {

	id := c.FormValue("id", "")
	token := c.FormValue("token", "")
	chunkStr := c.FormValue("chunk", "")
	if id == "" || token == "" || chunkStr == "" {
		return integration.InvalidRequest(c, "id, token and chunk are required")
	}

	chunk, err := strconv.Atoi(chunkStr)
	if err != nil {
		return integration.InvalidRequest(c, "invalid chunk")
	}

	transaction, ok := zapshare.GetTransaction(id)
	if !ok {
		return integration.InvalidRequest(c, "invalid id")
	}

	if transaction.Token != token {
		return integration.InvalidRequest(c, "invalid token")
	}

	return c.SendFile(transaction.ChunkFilePath(int64(chunk)))
}
