package magic_util

import (
	"log"
	"testing"

	"github.com/Liphium/station/neogate"
	"golang.org/x/net/websocket"
)

type GateConnection struct {
	conn *websocket.Conn
}

// Create a new gate connection (FOR TESTING ONLY)
func NewGateConnection(t *testing.T, token string) *GateConnection {
	ws, err := websocket.Dial(GatewayURL(), "", "http://localhost/")
	if err != nil {
		log.Fatalln("couldn't connect to gate:", err)
	}
	connection := &GateConnection{
		conn: ws,
	}

	_, err = ws.Write(Marshal(neogate.AuthPacket{
		Token:       token,
		Attachments: "",
	}))
	WebsocketError(t, err)

	ev := connection.ReadEvent(t)
	AssertEq(t, ev.Name, "ng_success")

	return connection
}

// Read a new event from the gate connection (FOR TESTING ONLY)
func (g *GateConnection) ReadEvent(t *testing.T) neogate.Event {
	// Set 1 second timeout for read operation
	//g.conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	buffer := make([]byte, 8*1024)
	n, err := g.conn.Read(buffer)
	WebsocketError(t, err)

	var event neogate.Event
	Unmarshal(buffer[:n], &event)
	return event
}

// Close the connection to the gate
func (g *GateConnection) Close(t *testing.T) {
	WebsocketError(t, g.conn.Close())
}
