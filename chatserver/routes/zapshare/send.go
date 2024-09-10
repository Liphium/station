package zapshare_routes

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
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

	// Get the file from the form
	file, err := c.FormFile("part")
	if err != nil {
		return integration.InvalidRequest(c, "no file")
	}

	// Make sure the file isn't too small
	if file.Size > zapshare.MaxChunkSize {
		return integration.InvalidRequest(c, "file too large")
	}

	// Make sure there aren't any file path things in the filename that could cause weird issues
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

	// Get file from the multipart header
	handle, err := file.Open()
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	defer handle.Close()

	// Extract the bytes from it
	bytes := make([]byte, file.Size)
	size, err := handle.Read(bytes)
	log.Println(size, "read")
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if the size read is the same as the actual file size
	if size != int(file.Size) {
		return integration.FailedRequest(c, localization.ErrorServer, errors.New("couldn't read the full chunk, maybe better logic is required?"))
	}

	// Add the chunk to the file parts cache
	transaction.FileParts.Store(chunk, &bytes)

	if err := transaction.PartUploaded(chunk); err != nil {
		return integration.InvalidRequest(c, "failed to update transaction")
	}

	return integration.SuccessfulRequest(c)
}
