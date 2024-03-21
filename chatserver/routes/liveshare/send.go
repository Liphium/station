package liveshare_routes

import (
	"fmt"
	"strings"

	"github.com/Liphium/station/chatserver/liveshare"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

func sendFilePart(c *fiber.Ctx) error {

	id := c.FormValue("id", "")
	token := c.FormValue("token", "")
	if id == "" || token == "" {
		return integration.InvalidRequest(c, "id and token are required")
	}

	transaction, valid := liveshare.GetTransaction(id)
	if !valid {
		return integration.InvalidRequest(c, "invalid transaction id")
	}

	if transaction.UploadToken != token {
		return integration.InvalidRequest(c, "invalid token")
	}

	file, err := c.FormFile("part")
	if err != nil {
		return integration.InvalidRequest(c, "no file")
	}

	if file.Size > liveshare.ChunkSize {
		return integration.InvalidRequest(c, "file too large")
	}

	if strings.Contains(file.Filename, "/") || strings.Contains(file.Filename, "\\") {
		return integration.InvalidRequest(c, "invalid filename")
	}

	fileName := fmt.Sprintf("chunk_%d", transaction.CurrentIndex)
	if file.Filename != fileName {
		return integration.InvalidRequest(c, "no chunk file name")
	}

	return c.SaveFile(file, transaction.VolumePath)
}
