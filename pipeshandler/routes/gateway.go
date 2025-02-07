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
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	router.Get("/", websocket.New(func(c *websocket.Conn) {
		ws(c, localNode, instance)
	}))
}

func ws(conn *websocket.Conn, local *pipes.LocalNode, instance *pipeshandler.Instance) {

	defer func() {
		util.PrintIfTesting("failed connection attempt")
		if err := recover(); err != nil {
			util.Log.Println("There was an error with a connection: ", err)
		}

		// Close the connection
		conn.Close()
	}()

	// Let the connection time out after 30 seconds
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))

	// Read the auth packet
	var authPacket struct {
		Token       string `json:"token"`
		Attachments string `json:"attachments"`
	}
	if err := conn.ReadJSON(&authPacket); err != nil {
		util.PrintIfTesting("closed connection: couldn't decode auth packet: ", err)
		return
	}

	// Check if the token is valid
	tk, ok := instance.CheckToken(authPacket.Token, local)
	if !ok {
		util.PrintIfTesting("closed connection: invalid auth token")
		return
	}

	// Make sure the session isn't already connected
	if instance.ExistsConnection(tk.Account, tk.Session) {
		util.PrintIfTesting("closed connection: already connected")
		return
	}

	// Ask the node if the connection should be accepted
	if instance.Config.TokenValidateHandler(tk, authPacket.Attachments) {
		util.PrintIfTesting("closed connection: token was rejected by service")
		return
	}

	// Make sure there is an infinite read timeout again (1 week should be enough)
	conn.SetReadDeadline(time.Now().Add(time.Hour * 24 * 7))

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

	if instance.Config.ClientConnectHandler(client, authPacket.Attachments) {
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

			// Disconnect the user on error
			OnError: func(err error) {

				// Remove the adapter
				local.RemoveAdapterWS(tk.Account)
			},
		})
	}

	if instance.Config.ClientEnterNetworkHandler(client, authPacket.Attachments) {
		return
	}

	for {
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
			Locale:     body["lc"].(string), // Parse the locale
			Node:       local,
			Instance:   instance,
		}) {
			instance.ReportClientError(client, "couldn't handle action", errors.New(body["action"].(string)))
			return
		}
	}
}
