package send

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/adapter"
	"github.com/bytedance/sonic"
)

// ClientUDP is a function that sends a WS packet to the client
func Client(id string, event pipes.Event) {

	msg, err := sonic.Marshal(event)
	if err != nil {
		return
	}

	adapter.ReceiveWeb(id, event, msg)
}

// ClientUDP sends a message to a client through UDP
func ClientUDP(id string, event pipes.Event) {

	msg, err := sonic.Marshal(event)
	if err != nil {
		return
	}

	adapter.ReceiveUDP(id, event, msg)
}
