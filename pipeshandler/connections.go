package pipeshandler

import (
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

func (instance *Instance) SendEvent(c *Client, event pipes.Event) error {

	msg, err := sonic.Marshal(event)
	if err != nil {
		return err
	}

	if c.Mutex == nil {
		c.Mutex = &sync.Mutex{}
	}

	c.Mutex.Lock()
	err = instance.SendMessage(c.Conn, c, msg)
	c.Mutex.Unlock()
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

	_, add := instance.connectionsCache.Get(getKey(client.ID, client.Session))
	instance.connectionsCache.Set(getKey(client.ID, client.Session), client, 1)

	if add {
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

func (instance *Instance) Remove(id string, session string) {
	client, valid := instance.Get(id, session)
	if valid {
		client.Conn.Close()
	}
	instance.connectionsCache.Del(getKey(id, session))
	instance.removeSession(id, session)
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

		instance.SendMessage(client.Conn, client, msg)
	}
}

func (instance *Instance) SendSession(id string, session string, msg []byte) bool {
	client, valid := instance.Get(id, session)
	if !valid {
		return false
	}

	instance.SendMessage(client.Conn, client, msg)
	return true
}

func (instance *Instance) SendMessage(conn *websocket.Conn, client *Client, msg []byte) error {

	msg, err := instance.Config.ClientEncodingMiddleware(client, instance, msg)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.BinaryMessage, msg)
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
