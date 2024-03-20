package liveshare_routes

import (
	"os"
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

	mainPath := os.Getenv("CN_LS_REPO") + "/" + transaction.Id
	return c.SaveFile(file, mainPath+"/"+file.Filename)
}
