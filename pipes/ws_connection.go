package pipes

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"nhooyr.io/websocket"
)

type AdoptionRequest struct {
	Token    string `json:"tk"`
	Adopting Node   `json:"adpt"`
}

func (local *LocalNode) setupWSStore() {
	local.nodeWSConnections = sync.Map{}
}

func (local *LocalNode) ConnectToNodeWS(node Node) error {

	// Marshal adoption request
	adoptionRq, err := sonic.Marshal(AdoptionRequest{
		Token:    node.Token,
		Adopting: local.ToNode(),
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
	local.nodeWSConnections.Store(node.ID, c)

	Log.Printf("[ws] Outgoing event stream to node %s connected.", node.ID)
	return nil
}

func (local *LocalNode) RemoveNodeWS(node string) {

	// Check if connection exists
	obj, ok := local.nodeWSConnections.Load(node)
	if !ok {
		return
	}
	connection := obj.(*websocket.Conn)

	// Close connection
	connection.Close(websocket.StatusNormalClosure, "Node disconnected")

	// Remove connection from map
	local.nodeWSConnections.Delete(node)

	Log.Printf("[ws] Outgoing event stream to node %s disconnected.", node)
}

func (local *LocalNode) ExistsNodeWS(node string) bool {

	// Check if connection exists
	_, ok := local.nodeWSConnections.Load(node)
	return ok
}

func (local *LocalNode) GetNodeWS(node string) *websocket.Conn {

	// Check if connection exists
	obj, ok := local.nodeWSConnections.Load(node)
	if !ok {
		return nil
	}
	connection := obj.(*websocket.Conn)

	return connection
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
func (local *LocalNode) IterateNodesWS(f func(key string, value *websocket.Conn) bool) {
	local.nodeWSConnections.Range(func(key, value any) bool {
		return f(key.(string), value.(*websocket.Conn))
	})
}

func (local *LocalNode) ReceiveWSAdoption(request string) (Node, error) {

	// Unmarshal
	var adoptionRq AdoptionRequest
	err := sonic.Unmarshal([]byte(request), &adoptionRq)
	if err != nil {
		return Node{}, err
	}

	// Check token
	if adoptionRq.Token != local.Token {
		return Node{}, errors.New("invalid token")
	}

	Log.Printf("[ws] Incoming event stream from node %s connected.", adoptionRq.Adopting.ID)
	local.AddNode(adoptionRq.Adopting)

	// Connect output stream (if not already connected)
	if !local.ExistsNodeWS(adoptionRq.Adopting.ID) {

		go func() {
			time.Sleep(2 * time.Second)

			Log.Println("Connecting to", adoptionRq.Adopting.ID)

			local.IterateNodesWS(func(id string, value *websocket.Conn) bool {
				Log.Println(id)
				return true
			})

			if err := local.ConnectToNodeWS(adoptionRq.Adopting); err != nil {
				Log.Println(err)
			}
		}()
	}

	return adoptionRq.Adopting, nil
}
