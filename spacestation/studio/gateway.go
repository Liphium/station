package studio

import (
	"github.com/pion/webrtc/v4"
)

// Start the gate for a specific connection
func (s *Studio) startGateway(c *Client, peer *webrtc.PeerConnection) error {

	// Create a new data channel for fast events (events over udp)
	ordered := false
	maxPacketLifetime := uint16(500)
	pipesChan, err := peer.CreateDataChannel("pipes", &webrtc.DataChannelInit{
		Ordered:           &ordered,
		MaxPacketLifeTime: &maxPacketLifetime,
	})
	if err != nil {
		return err
	}

	// Listen for keep alive messages from the client
	pipesChan.OnMessage(func(msg webrtc.DataChannelMessage) {
		// TODO: Handle fast pipe events
	})

	// Listen for new tracks
	peer.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		// TODO: Add the track
	})

	return nil
}
