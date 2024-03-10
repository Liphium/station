package send

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/receive"
	"github.com/bytedance/sonic"
)

const ProtocolWS = "ws"
const ProtocolUDP = "udp"

func Pipe(protocol string, message pipes.Message) error {

	if pipes.DebugLogs {
		pipes.Log.Printf("sent on [%s] %s: %s", protocol, message.Channel.Channel, message.Event.Name)
	}

	// Marshal message for sending to other nodes
	msg, err := sonic.Marshal(message)
	if err != nil {
		return err
	}

	// Send to receivers on current node
	receive.HandleMessage(protocol, message)

	if message.Local {
		return nil
	}

	switch message.Channel.Channel {
	case "conversation":
		return sendToConversation(protocol, message, msg)

	case "broadcast":
		return sendBroadcast(protocol, message, msg)

	case "p2p":
		return sendP2P(protocol, message, msg)
	}

	return nil
}
