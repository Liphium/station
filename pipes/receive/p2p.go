package receive

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/adapter"
	"github.com/Liphium/station/pipes/receive/processors"
)

func receiveP2P(protocol string, message pipes.Message) {

	// Process the message
	msg := processors.ProcessMarshal(&message, message.Channel.Target[0])
	if msg == nil {
		return
	}

	// Send to receiver
	switch protocol {
	case "ws":
		adapter.ReceiveWeb(message.Channel.Target[0], message.Event, msg)

	case "udp":
		adapter.ReceiveUDP(message.Channel.Target[0], message.Event, msg)
	}
}
