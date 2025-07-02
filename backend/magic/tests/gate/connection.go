package magic_gate

import (
	"log"
	"testing"
	"time"

	"github.com/Liphium/magic/mconfig"
	magic_accounts "github.com/Liphium/station/backend/magic/scripts/accounts"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
	"github.com/Liphium/station/neogate"
	"golang.org/x/net/websocket"
)

// Do not call this function anything with Test, it will cause errors
func MagicConnection(t *testing.T, p *mconfig.Plan) {
	magic_util.PrepareEnvironment(p)
	magic_util.PrepareDBTest()

	t.Run("basic gate flow", func(t *testing.T) {
		magic_accounts.TestAccount(p, "test1")
		_, tk := magic_accounts.GetTestToken(p, "test1")

		// Validate connections work
		conn := NewGateConnection(t, tk)
		defer conn.Close(t)
	})
}

type GateConnection struct {
	conn *websocket.Conn
}

// Create a new gate connection (FOR TESTING ONLY)
func NewGateConnection(t *testing.T, token string) *GateConnection {
	ws, err := websocket.Dial(magic_util.GatewayURL(), "", "http://localhost/")
	if err != nil {
		log.Fatalln("couldn't connect to gate:", err)
	}
	connection := &GateConnection{
		conn: ws,
	}

	_, err = ws.Write(magic_util.Marshal(neogate.AuthPacket{
		Token:       token,
		Attachments: "",
	}))
	magic_util.WebsocketError(t, err)

	ev := connection.ReadEvent(t)
	magic_util.AssertEq(t, ev.Name, "ng_success")

	return connection
}

// Read a new event from the gate connection (FOR TESTING ONLY)
func (g *GateConnection) ReadEvent(t *testing.T) neogate.Event {
	// Set 1 second timeout for read operation
	g.conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	buffer := make([]byte, 8*1024)
	n, err := g.conn.Read(buffer)
	magic_util.WebsocketError(t, err)

	var event neogate.Event
	magic_util.Unmarshal(buffer[:n], &event)
	return event
}

// Close the connection to the gate
func (g *GateConnection) Close(t *testing.T) {
	magic_util.WebsocketError(t, g.conn.Close())
}
