package pipeshandler

import (
	"fmt"
	"time"

	"github.com/Liphium/station/pipes"
	"github.com/bytedance/sonic"
	"github.com/dgraph-io/ristretto"
)

type Instance struct {
	Config           Config
	connectionsCache *ristretto.Cache // ID:Session -> Client
	sessionsCache    *ristretto.Cache // ID -> Session list
	routes           map[string]func(Context)
}

// ! If the functions aren't implemented pipesfiber will panic
type Config struct {
	ExpectedConnections int64
	SessionDuration     time.Duration // How long a session should stay alive
	Secret              []byte        // JWT secret (for authorization)

	// Node handlers
	NodeDisconnectHandler func(node pipes.Node)

	// Client handlers
	ClientDisconnectHandler   func(client *Client)
	TokenValidateHandler      func(claims *ConnectionTokenClaims, attachments string) bool // Called before the websocket connection is accepted, returns if the client should be disconnected (true = disconnect)
	ClientConnectHandler      func(client *Client, attachments string) bool                // Called after websocket connection is established, returns if the client should be disconnected (true = disconnect)
	ClientEnterNetworkHandler func(client *Client, attachments string) bool                // Called after pipes adapter is registered, returns if the client should be disconnected (true = disconnect)

	// Codec middleware
	ClientEncodingMiddleware func(client *Client, instance *Instance, message []byte) ([]byte, error)
	DecodingMiddleware       func(client *Client, instance *Instance, message []byte) (Message, error)

	// Error handler
	ErrorHandler func(err error)
}

// Message received from the client
type Message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// Default pipes-fiber encoding middleware (using JSON)
func DefaultClientEncodingMiddleware(client *Client, message []byte) ([]byte, error) {
	return message, nil
}

// Default pipes-fiber decoding middleware (using JSON)
func DefaultDecodingMiddleware(client *Client, bytes []byte) (Message, error) {
	var message Message
	if err := sonic.Unmarshal(bytes, &message); err != nil {
		return Message{}, err
	}
	return message, nil
}

func Setup(config Config) *Instance {
	instance := &Instance{
		Config: config,
	}
	instance.SetupConnectionsCache(config.ExpectedConnections)
	instance.routes = make(map[string]func(Context))
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
