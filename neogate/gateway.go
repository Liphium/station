package neogate

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// The packet needed for authenticating at the beginning of a neogate connection.
type AuthPacket struct {
	Token       string `json:"token"`
	Attachments string `json:"attachments"`
}

// Mount the neogate gateway using a fiber router.
func (instance *Instance) MountGateway(router fiber.Router) {

	// Inject a middleware to check if the request is a websocket upgrade request
	router.Use("/", func(c *fiber.Ctx) error {

		// Check if it is a websocket upgrade request
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}

		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	// Mount an endpoint for actually receiving the websocket connection
	router.Get("/", websocket.New(func(c *websocket.Conn) {
		ws(c, instance)
	}))
}

// Handles the websocket connection
func ws(conn *websocket.Conn, instance *Instance) {

	defer func() {
		if err := recover(); err != nil {
			Log.Println("There was an error with a connection: ", err)
			debug.PrintStack()
		}

		// Close the connection
		conn.Close()
	}()

	// Let the connection time out after 30 seconds
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))

	// Read the auth packet
	var authPacket AuthPacket
	if err := conn.ReadJSON(&authPacket); err != nil {
		Log.Println("closed connection: couldn't decode auth packet: ", err)
		return
	}

	// Check if the token is valid
	info, ok := instance.Config.CheckToken(authPacket.Token, authPacket.Attachments)
	if !ok {
		Log.Println("closed connection: invalid auth token")
		return
	}

	// Make sure the session isn't already connected
	if instance.ExistsConnection(info.Account, info.Session) {
		Log.Println("closed connection: already connected")
		return
	}

	// Make sure there is an infinite read timeout again (1 week should be enough)
	conn.SetReadDeadline(time.Now().Add(time.Hour * 24 * 7))

	client := instance.AddClient(info.ToClient(conn))
	defer func() {

		// Recover from a failure (in case of a cast issue maybe?)
		if err := recover(); err != nil {
			Log.Println("connection with", client.ID, "crashed cause of:", err)
		}

		// Get the client
		client, valid := instance.Get(info.Account, info.Session)
		if !valid {
			return
		}

		// Remove the connection from the cache
		instance.Config.ClientDisconnectHandler(client)
		instance.Remove(info.Account, info.Session)

		// Only remove adapter if all sessions are gone
		if len(instance.GetSessions(info.Account)) == 0 {
			instance.RemoveAdapter(info.Account)
		}
	}()

	if instance.Config.ClientConnectHandler(client, authPacket.Attachments) {
		return
	}

	// Add adapter for pipes (if this is the first session)
	if len(instance.GetSessions(info.Account)) == 1 {
		adapterName := instance.Config.ClientAdapterHandler(client)
		instance.Adapt(CreateAction{
			ID: adapterName,
			OnEvent: func(c *AdapterContext) error {
				if err := instance.SendToAccount(info.Account, c.Message); err != nil {
					instance.ReportClientError(client, "couldn't send received message", err)
					return err
				}
				return nil
			},

			// Disconnect the user on error
			OnError: func(err error) {
				instance.RemoveAdapter(adapterName)
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
			client, valid := instance.Get(info.Account, info.Session)
			if !valid {
				instance.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", info.Account, info.Session))
				return
			}

			instance.ReportClientError(client, "couldn't read message", err)
			return
		}

		// Get the client
		client, valid := instance.Get(info.Account, info.Session)
		if !valid {
			instance.ReportGeneralError("couldn't get client", fmt.Errorf("%s (%s)", info.Account, info.Session))
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
		if !instance.Handle(&Context{
			Client:     client,
			Data:       message,
			Action:     args[0],
			ResponseId: args[1],
			Locale:     body["lc"].(string), // Parse the locale
			Instance:   instance,
		}) {
			instance.ReportClientError(client, "couldn't handle action", errors.New(body["action"].(string)))
			return
		}
	}
}
