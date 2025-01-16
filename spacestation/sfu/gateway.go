package sfu

import "github.com/pion/webrtc/v4"

// Configuration
const keepAliveMessage = "liphium_spaces"

func manageConnection(room string, client string, peer *webrtc.PeerConnection) error {

	// Create a new data channel
	ordered := false
	maxPacketLifetime := uint16(1000)
	keepAliveChan, err := peer.CreateDataChannel("keepalive", &webrtc.DataChannelInit{
		Ordered:           &ordered,
		MaxPacketLifeTime: &maxPacketLifetime,
	})
	if err != nil {
		return err
	}

	// Listen for keep alive messages from the client
	keepAliveChan.OnMessage(func(msg webrtc.DataChannelMessage) {

		// Make sure the keep alive message is returning the exact same message
		if !msg.IsString || string(msg.Data) != keepAliveMessage {

		}
		updateKeepAlive(room, client)
	})

	// Start a goroutine that sends keep alive messages every 2 seconds
	go func() {
		if err := keepAliveChan.SendText(keepAliveMessage); err != nil {
			logger.Println("Couldn't send keep alive to", client, ": Ending connection.")
		}

	}()

	return nil
}
