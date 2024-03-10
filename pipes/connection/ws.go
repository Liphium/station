package connection

import (
	"context"

	"github.com/Liphium/station/pipes"
	"github.com/bytedance/sonic"
	"github.com/cornelk/hashmap"
	"nhooyr.io/websocket"
)

var nodeWSConnections = hashmap.New[string, *websocket.Conn]()

type AdoptionRequest struct {
	Token    string     `json:"tk"`
	Adopting pipes.Node `json:"adpt"`
}

func ConnectWS(node pipes.Node) error {

	// Marshal adoption request
	adoptionRq, err := sonic.Marshal(AdoptionRequest{
		Token:    node.Token,
		Adopting: pipes.CurrentNode,
	})
	if err != nil {
		return err
	}

	// Connect to node
	c, _, err := websocket.Dial(context.Background(), node.WS, &websocket.DialOptions{
		Subprotocols: []string{string(adoptionRq)},
	})

	if err != nil {
		return err
	}

	// Add connection to map
	nodeWSConnections.Insert(node.ID, c)

	pipes.Log.Printf("[ws] Outgoing event stream to node %s connected.", node.ID)
	return nil
}

func RemoveWS(node string) {

	// Check if connection exists
	connection, ok := nodeWSConnections.Get(node)
	if !ok {
		return
	}

	// Close connection
	connection.Close(websocket.StatusNormalClosure, "Node disconnected")

	// Remove connection from map
	nodeWSConnections.Del(node)

	pipes.Log.Printf("[ws] Outgoing event stream to node %s disconnected.", node)
}

func ExistsWS(node string) bool {

	// Check if connection exists
	_, ok := nodeWSConnections.Get(node)
	return ok
}

func GetWS(node string) *websocket.Conn {

	// Check if connection exists
	connection, ok := nodeWSConnections.Get(node)
	if !ok {
		return nil
	}

	return connection
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
func IterateWS(f func(key string, value *websocket.Conn) bool) {
	nodeWSConnections.Range(f)
}
