package liveshare_routes

import (
	"bufio"
	"fmt"

	"github.com/Liphium/station/chatserver/liveshare"
	"github.com/Liphium/station/main/integration"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

var Prefix = []byte("data: ")
var Suffix = []byte("\n\n")

// 6 cause "data", and 2 cause "\n\n"
const PacketSize = liveshare.ChunkSize + 6 + 2

func subscribeToLiveshare(c *fiber.Ctx) error {

	id := c.FormValue("id", "")
	token := c.FormValue("token", "")
	if id == "" || token == "" {
		return integration.InvalidRequest(c, "id and token are required")
	}

	receiverId, valid := liveshare.NewTransactionReceiver(id, token)
	if !valid {
		return integration.InvalidRequest(c, "Invalid id or token")

	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		fmt.Println("WRITER")

		firstPacket, err := sonic.Marshal(map[string]interface{}{
			"id": receiverId,
		})
		if err != nil {
			fmt.Printf("Error while marshalling: %v. Closing http connection.\n", err)
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", firstPacket)

		for {
			packet := make([]byte, PacketSize)
			copy(packet[:], Prefix)
			copy(packet[liveshare.ChunkSize:], Suffix)

			written, err := w.Write(packet)
			if err != nil {
				fmt.Printf("Error while writing: %v. Closing http connection.\n", err)
				break
			}
			fmt.Println(written)

			err = w.Flush()
			if err != nil {
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
				break
			}
		}
	}))

	return nil
}
