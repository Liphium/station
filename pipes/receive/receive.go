package receive

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/connection"
	"github.com/bytedance/sonic"
)

func ReceiveWS(bytes []byte) {

	// Decrypt
	var err error = nil
	bytes, err = connection.Decrypt(pipes.CurrentNode.ID, bytes)
	if err != nil {
		// TODO: Maybe report this?
		return
	}

	// Unmarshal
	var message pipes.Message
	err = sonic.Unmarshal(bytes, &message)
	if err != nil {
		return
	}

	// Handle message
	HandleMessage("ws", message)
}

func ReceiveUDP(bytes []byte) error {

	// Decrypt
	var err error = nil
	bytes, err = connection.Decrypt(pipes.CurrentNode.ID, bytes)
	if err != nil {
		return err
	}

	// Check for adoption request
	if bytes[0] == 'a' {

		// Adopt node
		return AdoptUDP(bytes)
	}

	// Unmarshal
	var message pipes.Message
	err = sonic.Unmarshal(bytes, &message)
	if err != nil {
		return err
	}

	// Handle message
	HandleMessage("udp", message)
	return nil
}

func HandleMessage(protocol string, message pipes.Message) {

	if pipes.DebugLogs {
		pipes.Log.Printf("received on [%s] %s: %s to %s", protocol, message.Channel.Channel, message.Event.Name, message.Channel.Target)
	}

	switch message.Channel.Channel {
	case pipes.ChannelBroadcast:
		receiveBroadcast(protocol, message)

	case pipes.ChannelConversation:
		receiveConversation(protocol, message)

	case pipes.ChannelP2P:
		receiveP2P(protocol, message)
	}
}
