package pipeshroutes

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Liphium/station/pipes/adapter"
	"github.com/Liphium/station/pipeshandler"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
	"github.com/Liphium/station/pipeshandler/wshandler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func gatewayRouter(router fiber.Router) {

	// Inject a middleware to check if the request is a websocket upgrade request
	router.Use("/", func(c *fiber.Ctx) error {

		// Check if it is a websocket upgrade request
		if websocket.IsWebSocketUpgrade(c) {

			// Check if the request has a token
			protocolSeperated := c.Get("Sec-WebSocket-Protocol")
			protocols := strings.Split(protocolSeperated, ", ")
			token := protocols[0]

			// Get attachments from the connection (passed to the node)
			attachments := ""
			if len(protocols) > 1 {
				attachments = protocols[1]
			}

			if len(token) == 0 {
				return c.SendStatus(fiber.StatusUnauthorized)
			}

			// Check if the token is valid
			tk, ok := pipeshandler.CheckToken(token)
			if !ok {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			// Make sure the session isn't already connected
			if pipeshandler.ExistsConnection(tk.Account, tk.Session) {
				return c.SendStatus(fiber.StatusConflict)
			}

			// Ask the node if the connection should be accepted
			if pipeshandler.CurrentConfig.TokenValidateHandler(tk, attachments) {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			// Set the token as a local variable
			c.Locals("ws", true)
			c.Locals("tk", tk)
			c.Locals("attached", attachments)
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	router.Get("/", websocket.New(ws))
}

func ws(conn *websocket.Conn) {
	tk := conn.Locals("tk").(*pipeshandler.ConnectionTokenClaims)

	client := pipeshandler.AddClient(tk.ToClient(conn, time.Now().Add(pipeshandler.CurrentConfig.SessionDuration)))
	defer func() {

		// Send callback to app
		client, valid := pipeshandler.Get(tk.Account, tk.Session)
		if !valid {
			return
		}
		adapter.RemoveWS(tk.Account)
		pipeshandler.CurrentConfig.ClientDisconnectHandler(client)

		// Remove the connection from the cache
		pipeshandler.Remove(tk.Account, tk.Session)
	}()

	if pipeshandler.CurrentConfig.ClientConnectHandler(client, conn.Locals("attached").(string)) {
		return
	}

	// Add adapter for pipes
	adapter.AdaptWS(adapter.Adapter{
		ID: tk.Account,
		Receive: func(c *adapter.Context) error {

			// Get the client
			client, valid := pipeshandler.Get(tk.Account, tk.Session)
			if !valid {
				pipeshandler.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", tk.Account, tk.Session))
				return errors.New("couldn't get client")
			}

			// Send message encoded with client encoding middleware
			msg, err := pipeshandler.CurrentConfig.ClientEncodingMiddleware(client, c.Message)
			if err != nil {
				pipeshandler.ReportClientError(client, "couldn't encode received message", err)
				return err
			}

			pipeshutil.Log.Println("sending "+c.Event.Name, "to", tk.Account)

			return conn.WriteMessage(websocket.BinaryMessage, msg)
		},
	})

	if pipeshandler.CurrentConfig.ClientEnterNetworkHandler(client, conn.Locals("attached").(string)) {
		return
	}

	for {

		// Read message as text
		_, msg, err := conn.ReadMessage()
		if err != nil {

			// Get the client for error reporting purposes
			client, valid := pipeshandler.Get(tk.Account, tk.Session)
			if !valid {
				pipeshandler.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", tk.Account, tk.Session))
				return
			}

			pipeshandler.ReportClientError(client, "couldn't read message", err)
			break
		}

		// Get the client
		client, valid := pipeshandler.Get(tk.Account, tk.Session)
		if !valid {
			pipeshandler.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", tk.Account, tk.Session))
			return
		}

		// Unmarshal the action
		message, err := pipeshandler.CurrentConfig.DecodingMiddleware(client, msg)
		if err != nil {
			pipeshandler.ReportClientError(client, "couldn't decode message", err)
			return
		}

		if client.IsExpired() {
			return
		}

		// Handle the action
		if !wshandler.Handle(wshandler.Message{
			Client: client,
			Data:   message.Data,
			Action: message.Action,
		}) {
			pipeshandler.ReportClientError(client, "couldn't handle action", errors.New(message.Action))
			return
		}
	}
}
