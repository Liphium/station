package zapshare_routes

import (
	"strconv"
	"strings"

	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

func sendFilePart(c *fiber.Ctx) error {

	id := c.FormValue("id", "")
	token := c.FormValue("token", "")
	if id == "" || token == "" {
		return integration.InvalidRequest(c, "id and token are required")
	}

	transaction, valid := zapshare.GetTransaction(id)
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

	if file.Size > zapshare.ChunkSize {
		return integration.InvalidRequest(c, "file too large")
	}

	if strings.Contains(file.Filename, "/") || strings.Contains(file.Filename, "\\") {
		return integration.InvalidRequest(c, "invalid filename")
	}

	args := strings.Split(file.Filename, "_")
	if len(args) != 2 {
		return integration.InvalidRequest(c, "invalid filename format")
	}
	chunkStr := args[1]
	chunk32, err := strconv.Atoi(chunkStr)
	if err != nil {
		return integration.InvalidRequest(c, "invalid chunk index")
	}
	chunk := int64(chunk32)

	if chunk > transaction.CurrentIndex+zapshare.ChunksAhead || chunk < transaction.CurrentIndex {
		return integration.InvalidRequest(c, "wrong chunk index")
	}

	if err := c.SaveFile(file, transaction.VolumePath+file.Filename); err != nil {
		return integration.InvalidRequest(c, "failed to save file")
	}

	if err := transaction.PartUploaded(chunk); err != nil {
		return integration.InvalidRequest(c, "failed to update transaction")
	}

	return nil
}
