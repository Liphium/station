package liveshare_routes

import (
	"bufio"
	"fmt"

	"github.com/Liphium/station/chatserver/liveshare"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func subscribeToLiveshare(c *fiber.Ctx) error {

	id := c.FormValue("id", "")
	token := c.FormValue("token", "")
	if id == "" || token == "" {
		return integration.InvalidRequest(c, "id and token are required")
	}

	receiver, valid := liveshare.NewTransactionReceiver(id, token)
	if !valid {
		return integration.InvalidRequest(c, "Invalid id or token")
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		util.Log.Println("Started subscription, waiting for packets...")

		// Send chunk start packet
		firstPacket, err := sonic.Marshal(map[string]interface{}{
			"id": receiver.ReceiverId,
		})
		if err != nil {
			util.Log.Println("Error while writing: ", err)
			return
		}
		_, err = w.Write(firstPacket)
		if err != nil {
			util.Log.Println("Error while writing: ", err)
			return
		}
		err = w.Flush()
		if err != nil {
			util.Log.Println("Error while flushing: ", err)
			return
		}

		for {
			packet := <-receiver.SendChannel

			// Send chunk data packet
			written, err := w.Write((*packet)[:])
			if err != nil {
				util.Log.Println("Error while writing: ", err)
				return
			}
			util.Log.Println("Wrote", written, "bytes to", receiver.ReceiverId)
			err = w.Flush()
			if err != nil {
				util.Log.Println("Error while flushing: ", err)
				return
			}

			// Send chunk end packet
			_, err = w.Write([]byte(fmt.Sprintf("\n\nc:%d\n", receiver.CurrentIndex)))
			if err != nil {
				util.Log.Println("Error while writing: ", err)
				return
			}
			err = w.Flush()
			if err != nil {
				util.Log.Println("Error while flushing: ", err)
				return
			}
		}
	}))

	return nil
}
