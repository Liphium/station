package pipeshandler

import (
	"fmt"
	"sync"
	"time"

	"github.com/Liphium/station/pipes"
	pipeshutil "github.com/Liphium/station/pipeshandler/util"
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

// Check the JWT token
func (instance *Instance) CheckToken(token string, local *pipes.LocalNode) (*ConnectionTokenClaims, bool) {

	// Check the jwt token
	jwtToken, err := jwt.ParseWithClaims(token, &ConnectionTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(instance.Config.Secret), nil
	})

	if err != nil {
		pipeshutil.Log.Println(err)
		return nil, false
	}

	// Check jwt claims
	if claims, ok := jwtToken.Claims.(*ConnectionTokenClaims); ok && jwtToken.Valid {

		// Validate the node id
		if claims.Node != local.ID {
			pipeshutil.Log.Println("invalid node")
			return nil, false
		}

		// Validate the expiration time
		if time.Now().After(time.Unix(claims.ExpiredUnixSec, 0)) {
			pipeshutil.Log.Println("invalid time")
			return nil, false
		}

		return claims, true
	}

	pipeshutil.Log.Println("invalid")
	return nil, false
}

// Generate a connection token for a node
func (instance *Instance) GenerateToken(account string, session string, extra string, node uint) (string, error) {

	// Create jwt token
	exp := time.Now().Add(time.Hour * 2)
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, ConnectionTokenClaims{
		Account:        account,
		ExpiredUnixSec: exp.Unix(),
		Session:        session,
		Extra:          extra,
		Node:           fmt.Sprintf("%d", node),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tk.SignedString([]byte(instance.Config.Secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
