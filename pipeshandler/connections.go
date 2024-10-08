package pipeshandler

import (
	"errors"
	"sync"
	"time"

	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
	"github.com/bytedance/sonic"
	"github.com/dgraph-io/ristretto"
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn    *websocket.Conn
	ID      string
	Session string
	End     time.Time
	Data    interface{}
	Mutex   *sync.Mutex
}

// Sends an event to the only ONE session of the connected account
func (instance *Instance) SendEventToOne(c *Client, event pipes.Event) error {

	msg, err := sonic.Marshal(event)
	if err != nil {
		return err
	}

	err = instance.SendMessage(c, msg)
	return err
}

func (c *Client) IsExpired() bool {
	return c.End.Before(time.Now())
}

func (instance *Instance) SetupConnectionsCache(expected int64) {

	var err error
	instance.connectionsCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: expected * 10, // pass in expected items
		MaxCost:     1 << 30,       // maximum cost of cache is 1GB
		BufferItems: 64,            // Some random number, check docs
	})

	if err != nil {
		panic(err)
	}

	instance.sessionsCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: expected * 10, // pass in expected items
		MaxCost:     1 << 30,       // maximum cost of cache is 1GB
		BufferItems: 64,            // Some random number, check docs
	})

	if err != nil {
		panic(err)
	}

}

func getKey(id string, session string) string {
	return id + ":" + session
}

func (instance *Instance) AddClient(client Client) *Client {

	// Add the session
	_, valid := instance.connectionsCache.Get(getKey(client.ID, client.Session))
	instance.connectionsCache.Set(getKey(client.ID, client.Session), client, 1)
	instance.connectionsCache.Wait()

	// If the session is not yet added, make sure to add it to the list
	if !valid {
		instance.addSession(client.ID, client.Session)
	}

	return &client
}

func (instance *Instance) UpdateClient(client *Client) {
	instance.connectionsCache.Set(getKey(client.ID, client.Session), *client, 1)
	instance.connectionsCache.Wait()
}

func (instance *Instance) GetSessions(id string) []string {
	sessions, valid := instance.sessionsCache.Get(id)
	if valid {
		return sessions.([]string)
	}

	return []string{}
}

func (instance *Instance) addSession(id string, session string) {

	sessions, valid := instance.sessionsCache.Get(id)
	if valid {
		instance.sessionsCache.Set(id, append(sessions.([]string), session), 1)
	} else {
		instance.sessionsCache.Set(id, []string{session}, 1)
	}
	instance.sessionsCache.Wait()
}

func (instance *Instance) removeSession(id string, session string) {

	sessions, valid := instance.sessionsCache.Get(id)
	if valid {

		if len(sessions.([]string)) == 1 {
			instance.sessionsCache.Del(id)
			return
		}

		instance.sessionsCache.Set(id, pipeshutil.RemoveString(sessions.([]string), session), 1)
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
	instance.connectionsCache.Del(getKey(id, session))
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

func (instance *Instance) Send(id string, msg []byte) {
	sessions, ok := instance.sessionsCache.Get(id)

	if !ok {
		return
	}

	for _, session := range sessions.([]string) {
		client, valid := instance.Get(id, session)
		if !valid {
			continue
		}

		instance.SendMessage(client, msg)
	}
}

func (instance *Instance) SendSession(id string, session string, msg []byte) bool {
	client, valid := instance.Get(id, session)
	if !valid {
		return false
	}

	instance.SendMessage(client, msg)
	return true
}

func (instance *Instance) SendMessage(client *Client, msg []byte) error {

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
	_, ok := instance.connectionsCache.Get(getKey(id, session))
	if !ok {
		return false
	}

	return ok
}

func (instance *Instance) Get(id string, session string) (*Client, bool) {
	client, valid := instance.connectionsCache.Get(getKey(id, session))
	if !valid {
		return &Client{}, false
	}

	cl := client.(Client)
	return &cl, true
}

func (instance *Instance) GetConnections(id string) int {
	clients, ok := instance.sessionsCache.Get(id)
	if !ok {
		return 0
	}

	return len(clients.([]string))
}
