package sfu

import "github.com/pion/webrtc/v4"

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
		// TODO: Handle somehow
	})

	// TODO: Handle connection and disconnection logic

	return nil
}
