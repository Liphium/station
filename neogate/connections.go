package neogate

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn    *websocket.Conn
	ID      string
	Session string
	Extra   interface{}
	End     time.Time
	Data    interface{}
	Mutex   *sync.Mutex
}

// Sends an event to the client
func (instance *Instance) SendEventToClient(c *Client, event Event) error {
	msg, err := sonic.Marshal(event)
	if err != nil {
		return err
	}

	err = instance.SendToClient(c, msg)
	return err
}

func (c *Client) IsExpired() bool {
	return c.End.Before(time.Now())
}

func getKey(id string, session string) string {
	return id + ":" + session
}

func (instance *Instance) AddClient(client Client) *Client {

	// Add the session
	_, valid := instance.connectionsCache.Load(getKey(client.ID, client.Session))
	instance.connectionsCache.Store(getKey(client.ID, client.Session), client)

	// If the session is not yet added, make sure to add it to the list
	if !valid {
		instance.addSession(client.ID, client.Session)
	}

	return &client
}

func (instance *Instance) UpdateClient(client *Client) {
	instance.connectionsCache.Store(getKey(client.ID, client.Session), *client)
}

func (instance *Instance) GetSessions(id string) []string {
	sessions, valid := instance.sessionsCache.Load(id)
	if valid {
		return sessions.([]string)
	}

	return []string{}
}

func (instance *Instance) addSession(id string, session string) {

	sessions, valid := instance.sessionsCache.Load(id)
	if valid {
		instance.sessionsCache.Store(id, append(sessions.([]string), session))
	} else {
		instance.sessionsCache.Store(id, []string{session})
	}
}

func (instance *Instance) removeSession(id string, session string) {

	sessions, valid := instance.sessionsCache.Load(id)
	if valid {

		if len(sessions.([]string)) == 1 {
			instance.sessionsCache.Delete(id)
			return
		}

		instance.sessionsCache.Store(id, slices.DeleteFunc(sessions.([]string), func(s string) bool {
			return s == session
		}))
	}
}

// Remove a session from the account (DOES NOT DISCONNECT, there is an extra method for that)
func (instance *Instance) Remove(id string, session string) {
	client, valid := instance.Get(id, session)
	if valid {
		err := client.Conn.Close()
		if err != nil {
			instance.ReportGeneralError("couldn't disconnect client", err)
		}
	} else {
		instance.ReportGeneralError("client "+id+" doesn't exist", errors.New("couldn't delete"))
	}
	instance.connectionsCache.Delete(getKey(id, session))
	instance.removeSession(id, session)
}

// Disconnect a client from the network
func (instance *Instance) Disconnect(id string, session string) {

	// Get the client
	client, valid := instance.Get(id, session)
	if !valid {
		return
	}

	// This is a little weird for disconnecting, but it works, so I'm not complaining
	client.Conn.SetReadDeadline(time.Now().Add(time.Microsecond * 1))
	client.Conn.Close()
}

// Send bytes to an account id
func (instance *Instance) SendToAccount(id string, msg []byte) error {
	sessions, ok := instance.sessionsCache.Load(id)
	if !ok {
		return errors.New("no sessions found")
	}

	for _, session := range sessions.([]string) {
		client, valid := instance.Get(id, session)
		if !valid {
			continue
		}

		if err := instance.SendToClient(client, msg); err != nil {
			return err
		}
	}
	return nil
}

func (instance *Instance) SendToSession(id string, session string, msg []byte) bool {
	client, valid := instance.Get(id, session)
	if !valid {
		return false
	}

	instance.SendToClient(client, msg)
	return true
}

func (instance *Instance) SendToClient(client *Client, msg []byte) error {

	msg, err := instance.Config.ClientEncodingMiddleware(client, instance, msg)
	if err != nil {
		return err
	}

	// Make sure there are no concurrent writes
	if client.Mutex == nil {
		client.Mutex = &sync.Mutex{}
	}

	// Lock and unlock mutex after writing
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	return client.Conn.WriteMessage(websocket.BinaryMessage, msg)
}

func (instance *Instance) ExistsConnection(id string, session string) bool {
	_, ok := instance.connectionsCache.Load(getKey(id, session))
	if !ok {
		return false
	}

	return ok
}

func (instance *Instance) Get(id string, session string) (*Client, bool) {
	client, valid := instance.connectionsCache.Load(getKey(id, session))
	if !valid {
		return &Client{}, false
	}

	cl := client.(Client)
	return &cl, true
}

func (instance *Instance) GetConnections(id string) int {
	clients, ok := instance.sessionsCache.Load(id)
	if !ok {
		return 0
	}

	return len(clients.([]string))
}
