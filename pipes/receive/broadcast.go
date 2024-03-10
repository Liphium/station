package receive

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/adapter"
	"github.com/Liphium/station/pipes/receive/processors"
)

func receiveBroadcast(protocol string, message pipes.Message) {

	if message.Event.Name == "ping" {
		pipes.Log.Println("Received ping from node", message.Event.Data["node"])
	}

	// Send to all receivers
	for _, tg := range message.Channel.Target {

		// Process the event
		msg := processors.ProcessMarshal(&message, tg)
		if msg == nil {
			continue
		}

		// Send to correct adapter
		switch protocol {
		case "ws":
			adapter.ReceiveWeb(tg, message.Event, msg)

		case "udp":
			adapter.ReceiveUDP(tg, message.Event, msg)
		}
	}
}
