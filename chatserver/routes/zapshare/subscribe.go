package zapshare_routes

import (
	"bufio"
	"fmt"

	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func subscribeToLiveshare(c *fiber.Ctx) error {

	id := c.FormValue("id", "")
	token := c.FormValue("token", "")
	if id == "" || token == "" {
		return integration.InvalidRequest(c, "id and token are required")
	}

	receiver, valid := zapshare.NewTransactionReceiver(id, token)
	if !valid {
		return integration.InvalidRequest(c, "Invalid id or token")
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		util.Log.Println("Started subscription, waiting for packets...")
		defer func() {
			if err := recover(); err != nil {
				util.Log.Println("Recovered from panic: ", err)
			}
			util.Log.Println("Cancelling zapshare session")
			zapshare.CancelTransaction(id)
		}()

		// Send chunk start packet
		_, err := fmt.Fprintf(w, "data: %s\n\n", receiver.ReceiverId)
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

			if packet == -1 {
				util.Log.Println("Subscription ended")
				return
			}

			// Send chunk data packet
			written, err := fmt.Fprintf(w, "data: %d\n\n", packet)
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
		}
	}))

	return nil
}
