package pipeshroutes

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func gatewayRouter(router fiber.Router, localNode *pipes.LocalNode, instance *pipeshandler.Instance) {

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
			tk, ok := instance.CheckToken(token, localNode)
			if !ok {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			// Make sure the session isn't already connected
			if instance.ExistsConnection(tk.Account, tk.Session) {
				return c.SendStatus(fiber.StatusConflict)
			}

			// Ask the node if the connection should be accepted
			if instance.Config.TokenValidateHandler(tk, attachments) {
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

	router.Get("/", websocket.New(func(c *websocket.Conn) {
		ws(c, localNode, instance)
	}))
}

func ws(conn *websocket.Conn, local *pipes.LocalNode, instance *pipeshandler.Instance) {
	tk := conn.Locals("tk").(*pipeshandler.ConnectionTokenClaims)

	client := instance.AddClient(tk.ToClient(conn, time.Now().Add(instance.Config.SessionDuration)))
	defer func() {

		// Recover from a failure (in case of a cast issue maybe?)
		if err := recover(); err != nil {
			util.Log.Println("connection with", client.ID, "crashed cause of:", err)
		}

		// Get the client
		client, valid := instance.Get(tk.Account, tk.Session)
		if !valid {
			return
		}

		// Remove the connection from the cache
		instance.Config.ClientDisconnectHandler(client)
		instance.Remove(tk.Account, tk.Session)

		// Only remove adapter if all sessions are gone
		if len(instance.GetSessions(tk.Account)) == 0 {
			local.RemoveAdapterWS(tk.Account)
		}
	}()

	if instance.Config.ClientConnectHandler(client, conn.Locals("attached").(string)) {
		return
	}

	// Add adapter for pipes (if this is the first session)
	if len(instance.GetSessions(tk.Account)) == 1 {
		local.AdaptWS(pipes.Adapter{
			ID: tk.Account,
			Receive: func(c *pipes.Context) error {
				for _, session := range instance.GetSessions(tk.Account) {
					log.Println("SENDING TO SESSION", session)
					// Get the client
					client, valid := instance.Get(tk.Account, session)
					if !valid {
						instance.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", tk.Account, session))
						return errors.New("couldn't get client")
					}

					// Send message encoded with client encoding middleware
					pipeshutil.Log.Println("sending "+c.Event.Name, "to", tk.Account)
					if err := instance.SendMessage(client, c.Message); err != nil {
						instance.ReportClientError(client, "couldn't send received message", err)
						return err
					}
				}

				return nil
			},
		})
	}

	if instance.Config.ClientEnterNetworkHandler(client, conn.Locals("attached").(string)) {
		return
	}

	for {
		log.Println("reading..")

		_, msg, err := conn.ReadMessage()
		if err != nil {

			// Get the client for error reporting purposes
			client, valid := instance.Get(tk.Account, tk.Session)
			if !valid {
				instance.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", tk.Account, tk.Session))
				return
			}

			instance.ReportClientError(client, "couldn't read message", err)
			return
		}

		// Get the client
		client, valid := instance.Get(tk.Account, tk.Session)
		if !valid {
			instance.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", tk.Account, tk.Session))
			return
		}

		// Decode the message
		message, err := instance.Config.DecodingMiddleware(client, instance, msg)
		if err != nil {
			instance.ReportClientError(client, "couldn't decode message", err)
			return
		}

		if client.IsExpired() {
			return
		}

		// Unmarshal the message to extract a few things
		var body map[string]interface{}
		if err := sonic.Unmarshal(message, &body); err != nil {
			return
		}

		// Extract the response id from the message
		args := strings.Split(body["action"].(string), ":")
		if len(args) != 2 {
			return
		}

		// Handle the action
		if !instance.Handle(&pipeshandler.Context{
			Client:     client,
			Data:       message,
			Action:     args[0],
			ResponseId: args[1],
			Node:       local,
			Instance:   instance,
		}) {
			instance.ReportClientError(client, "couldn't handle action", errors.New(body["action"].(string)))
			return
		}
	}
}
