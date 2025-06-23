package neogate

import (
	"fmt"
	"log"
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/dgraph-io/ristretto"
	"github.com/gofiber/websocket/v2"
)

var Log = log.New(log.Writer(), "neogate ", log.Flags())

type Instance struct {
	Config           Config
	connectionsCache *ristretto.Cache // ID:Session -> Client
	sessionsCache    *ristretto.Cache // ID -> Session list
	adapters         *sync.Map        // ID -> Adapter
	routes           map[string]func(*Context) pipes.Event
}

type ClientInfo struct {
	Account string      // Identifier of the account
	Session string      // Identifier of the session
	Extra   interface{} // Extra data that's up to you to fill
}

// Convert the client information to a client that can be used by neogate.
func (info ClientInfo) ToClient(conn *websocket.Conn) Client {
	return Client{
		Conn:    conn,
		ID:      info.Account,
		Session: info.Session,
		Extra:   info.Extra,
		Mutex:   &sync.Mutex{},
	}
}

// ! If the functions aren't implemented pipesfiber will panic
type Config struct {
	ExpectedConnections int64
	Secret              []byte // JWT secret (for authorization)

	// Called when a client attempts to connection using a token. Return true if the token is valid. MUST BE SPECIFIED.
	CheckToken func(token string, attachments string) (ClientInfo, bool)

	// Client handlers
	ClientDisconnectHandler   func(client *Client)
	ClientConnectHandler      func(client *Client, attachments string) bool // Called after websocket connection is established, returns if the client should be disconnected (true = disconnect)
	ClientEnterNetworkHandler func(client *Client, attachments string) bool // Called after pipes adapter is registered, returns if the client should be disconnected (true = disconnect)

	// Codec middleware
	ClientEncodingMiddleware func(client *Client, instance *Instance, message []byte) ([]byte, error)
	DecodingMiddleware       func(client *Client, instance *Instance, message []byte) ([]byte, error)

	// Error handler
	ErrorHandler func(err error)
}

// Message received from the client
type Message[T any] struct {
	Action string `json:"action"`
	Data   T      `json:"data"`
}

// Default pipes-fiber encoding middleware (using JSON)
func DefaultClientEncodingMiddleware(client *Client, message []byte) ([]byte, error) {
	return message, nil
}

// Default pipes-fiber decoding middleware (using JSON)
func DefaultDecodingMiddleware(client *Client, bytes []byte) ([]byte, error) {
	return bytes, nil
}

// Setup neogate using the config. Use the returned *Instance for interfacing with neogate.
func Setup(config Config) *Instance {
	instance := &Instance{
		Config:   config,
		adapters: &sync.Map{},
		routes:   make(map[string]func(*Context) pipes.Event),
	}
	instance.SetupConnectionsCache(config.ExpectedConnections)
	return instance
}

func (instance *Instance) ReportGeneralError(context string, err error) {
	if instance.Config.ErrorHandler == nil {
		return
	}

	instance.Config.ErrorHandler(fmt.Errorf("general: %s: %s", context, err.Error()))
}

func (instance *Instance) ReportClientError(client *Client, context string, err error) {
	if instance.Config.ErrorHandler == nil {
		return
	}

	instance.Config.ErrorHandler(fmt.Errorf("client %s: %s: %s", client.ID, context, err.Error()))
}
