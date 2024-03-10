package send

import (
	"context"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/connection"
	"nhooyr.io/websocket"
)

func sendToConversation(protocol string, message pipes.Message, msg []byte) error {

	for _, node := range message.Channel.Nodes {
		if node == pipes.CurrentNode.ID {
			continue
		}

		// Encrypt message for node
		encryptedMsg, err := connection.Encrypt(node, msg)
		if err != nil {
			return err
		}

		switch protocol {
		case "ws":
			connection.GetWS(node).Write(context.Background(), websocket.MessageText, encryptedMsg)

		case "udp":
			connection.GetUDP(node).Write(encryptedMsg)
		}
	}

	return nil
}
