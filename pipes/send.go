package pipes

import (
	"context"

	"github.com/Liphium/station/pipes/util"
	"github.com/bytedance/sonic"
	"nhooyr.io/websocket"
)

const ProtocolWS = "ws"
const ProtocolUDP = "udp"

func (node *LocalNode) Pipe(protocol string, message Message) error {

	if DebugLogs {
		Log.Printf("sent on [%s] %s: %s", protocol, message.Channel.Channel, message.Event.Name)
	}

	// Marshal message for sending to other nodes
	msg, err := sonic.Marshal(message)
	if err != nil {
		return err
	}

	// Send to receivers on current node
	//HandleMessage(protocol, message)

	if message.Local {
		return nil
	}

	switch message.Channel.Channel {
	case "conversation":
		return node.sendToConversation(protocol, message, msg)

	case "broadcast":
		return node.sendBroadcast(protocol, message, msg)

	case "p2p":
		return node.sendP2P(protocol, message, msg)
	}

	return nil
}

func (local *LocalNode) sendBroadcast(protocol string, message Message, msg []byte) error {

	// Send to other nodes
	var mainErr error = nil
	switch protocol {
	case "ws":
		local.IterateNodesWS(func(id string, node *websocket.Conn) bool {

			// Encrypt message for node
			encryptedMsg, err := local.Encrypt(id, msg)
			mainErr = err
			if err != nil {
				return false
			}

			node.Write(context.Background(), websocket.MessageText, encryptedMsg)
			return true
		})

	}

	return mainErr
}

func (local *LocalNode) sendToConversation(protocol string, message Message, msg []byte) error {

	for _, node := range message.Channel.Nodes {
		if node == local.ID {
			continue
		}

		// Encrypt message for node
		encryptedMsg, err := local.Encrypt(node, msg)
		if err != nil {
			return err
		}

		switch protocol {
		case "ws":
			local.GetNodeWS(node).Write(context.Background(), websocket.MessageText, encryptedMsg)
		}
	}

	return nil
}

func (local *LocalNode) sendP2P(protocol string, message Message, msg []byte) error {

	// Check if receiver is on this node
	if message.Channel.Target[0] == local.ID {
		switch protocol {
		case "ws":
			local.AdapterReceiveWeb(message.Channel.Target[1], message.Event, msg)
		}
		return nil
	}

	// Encrypt message for node
	encryptedMsg, err := local.Encrypt(message.Channel.Target[1], msg)
	if err != nil {
		return err
	}

	// Send to correct node
	switch protocol {
	case "ws":
		local.GetNodeWS(message.Channel.Target[1]).Write(context.Background(), websocket.MessageText, encryptedMsg)
	}

	return nil
}

// SendClient is a function that sends a WS packet to the client
func (local *LocalNode) SendClient(id string, event Event) {

	msg, err := sonic.Marshal(event)
	if err != nil {
		return
	}

	local.AdapterReceiveWeb(id, event, msg)
}

func (local *LocalNode) Socketless(nodeEntity Node, message Message) error {

	if DebugLogs {
		Log.Printf("sent on [socketless] %s: %s", message.Channel.Channel, message.Event.Name)
	}

	if nodeEntity.ID == local.ID {

		//local.HandleMessage("ws", message)
		return nil
	}

	err := util.PostRaw(nodeEntity.SL, map[string]interface{}{
		"token":   nodeEntity.Token,
		"message": message,
	})

	if err != nil {
		return err
	}

	return nil
}
