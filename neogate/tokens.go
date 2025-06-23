package neogate

import (
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Connection token struct
type ConnectionTokenClaims struct {
	Account        string `json:"acc"`             // Account id of the connecting client
	ExpiredUnixSec int64  `json:"e_u"`             // Expiration time in unix seconds
	Session        string `json:"ses"`             // Session id of the connecting client
	Node           string `json:"node"`            // Node id of the node the client is connecting to
	Extra          string `json:"extra,omitempty"` // Extra arguments for the connection

	jwt.RegisteredClaims
}

func (tk ConnectionTokenClaims) ToClient(conn *websocket.Conn, end time.Time) Client {
	return Client{
		Conn:    conn,
		ID:      tk.Account,
		Session: tk.Session,
		End:     end,
		Mutex:   &sync.Mutex{},
	}
}
