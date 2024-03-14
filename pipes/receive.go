package pipes

import (
	"github.com/bytedance/sonic"
)

func (local *LocalNode) ReceiveWS(bytes []byte) {

	// Decrypt
	var err error = nil
	bytes, err = local.Decrypt(local.ID, bytes)
	if err != nil {
		// TODO: Maybe report this?
		return
	}

	// Unmarshal
	var message Message
	err = sonic.Unmarshal(bytes, &message)
	if err != nil {
		return
	}

	// Handle message
	local.HandleMessage(ProtocolWS, message)
}

func (local *LocalNode) HandleMessage(protocol string, message Message) {

	if DebugLogs {
		Log.Printf("received on [%s] %s: %s to %s", protocol, message.Channel.Channel, message.Event.Name, message.Channel.Target)
	}

	switch message.Channel.Channel {
	case ChannelBroadcast:
		local.receiveBroadcast(protocol, message)

	case ChannelConversation:
		local.receiveConversation(protocol, message)

	case ChannelP2P:
		local.receiveP2P(protocol, message)
	}
}

func (local *LocalNode) receiveBroadcast(protocol string, message Message) {

	if message.Event.Name == "ping" {
		Log.Println("Received ping from node", message.Event.Data["node"])
	}

	// Send to all receivers
	for _, tg := range message.Channel.Target {

		// Process the event
		msg := local.ProcessMarshal(&message, tg)
		if msg == nil {
			continue
		}

		// Send to correct adapter
		switch protocol {
		case "ws":
			local.AdapterReceiveWeb(tg, message.Event, msg)
		}
	}
}

func (local *LocalNode) receiveConversation(protocol string, message Message) {

	// Send to receivers
	for _, member := range message.Channel.Target {

		// Process the message
		msg := local.ProcessMarshal(&message, member)
		if msg == nil {
			continue
		}

		// Send to correct adapter
		switch protocol {
		case "ws":
			local.AdapterReceiveWeb(member, message.Event, msg)
		}
	}
}

func (local *LocalNode) receiveP2P(protocol string, message Message) {

	// Process the message
	msg := local.ProcessMarshal(&message, message.Channel.Target[0])
	if msg == nil {
		return
	}

	// Send to receiver
	switch protocol {
	case "ws":
		local.AdapterReceiveWeb(message.Channel.Target[0], message.Event, msg)

	}
}
