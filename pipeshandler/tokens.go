package pipeshandler

import (
	"log"
	"sync"
	"time"

	"github.com/Liphium/station/pipes"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Connection token struct
type ConnectionTokenClaims struct {
	Account        string `json:"acc"`  // Account id of the connecting client
	ExpiredUnixSec int64  `json:"e_u"`  // Expiration time in unix seconds
	Session        string `json:"ses"`  // Session id of the connecting client
	Node           string `json:"node"` // Node id of the node the client is connecting to

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
func CheckToken(token string) (*ConnectionTokenClaims, bool) {

	// Check the jwt token
	jwtToken, err := jwt.ParseWithClaims(token, &ConnectionTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(CurrentConfig.Secret), nil
	})

	if err != nil {
		log.Println(err)
		return nil, false
	}

	// Check jwt claims
	if claims, ok := jwtToken.Claims.(*ConnectionTokenClaims); ok && jwtToken.Valid {

		// Validate the node id
		if claims.Node != pipes.CurrentNode.ID {
			log.Println("invalid node")
			return nil, false
		}

		// Validate the expiration time
		if time.Now().After(time.Unix(claims.ExpiredUnixSec, 0)) {
			log.Println("invalid time")
			return nil, false
		}

		return claims, true
	}

	log.Println("invalid")
	return nil, false
}
