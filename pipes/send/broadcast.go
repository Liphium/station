package send

import (
	"context"
	"net"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/connection"
	"nhooyr.io/websocket"
)

func sendBroadcast(protocol string, message pipes.Message, msg []byte) error {

	// Send to other nodes
	var mainErr error = nil
	switch protocol {
	case "ws":
		connection.IterateWS(func(id string, node *websocket.Conn) bool {

			// Encrypt message for node
			encryptedMsg, err := connection.Encrypt(id, msg)
			mainErr = err
			if err != nil {
				return false
			}

			node.Write(context.Background(), websocket.MessageText, encryptedMsg)
			return true
		})

	case "udp":
		connection.IterateUDP(func(id string, node *net.UDPConn) bool {

			// Encrypt message for node
			encryptedMsg, err := connection.Encrypt(id, msg)
			mainErr = err
			if err != nil {
				return false
			}

			node.Write(encryptedMsg)
			return true
		})
	}

	return mainErr
}
