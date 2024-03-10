package connection

import (
	"net"

	"github.com/Liphium/station/pipes"
	"github.com/bytedance/sonic"
	"github.com/cornelk/hashmap"
)

var nodeUDPConnections = hashmap.New[string, *net.UDPConn]()

/* Eventually implement custom UDP protocol
type AdoptionRequest struct {
	Token    string     `json:"tk"`
	Adopting pipes.Node `json:"adpt"`
}
*/

var GeneralPrefix []byte = nil

func ConnectUDP(node pipes.Node) error {

	// Marshal current node
	adoptionRq, err := sonic.Marshal(AdoptionRequest{
		Token:    node.Token,
		Adopting: pipes.CurrentNode,
	})
	if err != nil {
		return err
	}

	// Add prefix
	adoptionRq = append([]byte("a:"), adoptionRq...)

	// Encrypt
	adoptionRq, err = Encrypt(node.ID, adoptionRq)
	if err != nil {
		return err
	}

	// Resolve udp address
	udpAddr, err := net.ResolveUDPAddr("udp", node.UDP)
	if err != nil {
		return err
	}

	// Connect to node
	c, err := net.DialUDP("udp", nil, udpAddr)

	if err != nil {
		c.Close()
		return err
	}

	// Send adoption request
	_, err = c.Write(adoptionRq)
	if err != nil {
		c.Close()
		return err
	}

	// Add connection to map
	nodeUDPConnections.Insert(node.ID, c)

	pipes.Log.Printf("[udp] Outgoing event stream to node %s connected.", node.ID)
	return nil
}

func RemoveUDP(node string) {

	// Check if connection exists
	connection, ok := nodeUDPConnections.Get(node)
	if !ok {
		return
	}

	// Close connection
	connection.Close()

	// Remove connection from map
	nodeUDPConnections.Del(node)
}

func ExistsUDP(node string) bool {

	// Check if connection exists
	_, ok := nodeUDPConnections.Get(node)
	return ok
}

func GetUDP(node string) *net.UDPConn {

	// Check if connection exists
	connection, ok := nodeUDPConnections.Get(node)
	if !ok {
		return nil
	}

	return connection
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
func IterateUDP(f func(key string, value *net.UDPConn) bool) {
	nodeUDPConnections.Range(f)
}
