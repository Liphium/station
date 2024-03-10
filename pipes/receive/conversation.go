package receive

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/adapter"
	"github.com/Liphium/station/pipes/receive/processors"
)

func receiveConversation(protocol string, message pipes.Message) {

	// Send to receivers
	for _, member := range message.Channel.Target {

		// Process the message
		msg := processors.ProcessMarshal(&message, member)
		if msg == nil {
			continue
		}

		// Send to correct adapter
		switch protocol {
		case "ws":
			adapter.ReceiveWeb(member, message.Event, msg)

		case "udp":
			adapter.ReceiveUDP(member, message.Event, msg)
		}
	}
}
