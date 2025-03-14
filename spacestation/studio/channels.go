package studio

import (
	"errors"
	"io"
	"sync"

	"github.com/pion/webrtc/v4"
)

type Channel struct {
	id            string                      // Id of the channel (read-only)
	track         *Track                      // Track the channel is attached to (read-only)
	remoteTrack   *webrtc.TrackRemote         // read-only
	localTrack    *webrtc.TrackLocalStaticRTP // read-only
	subscriptions *sync.Map                   // Client id -> *Subscription (read-only)
}

// Creates the tracks required for operating the channel (must be called before starting the sender)
func (c *Channel) Init() error {

	// Create a new track for sending
	track, err := webrtc.NewTrackLocalStaticRTP(c.remoteTrack.Codec().RTPCodecCapability, c.track.id, c.id)
	if err != nil {
		return err
	}

	// Set the track
	c.localTrack = track

	return nil
}

// Start the sender that forwards the packets
func (c *Channel) startSender() {
	for {
		// Read RTP packets being sent on the channel
		packet, _, readErr := c.remoteTrack.ReadRTP()
		if readErr != nil {
			logger.Println("Couldn't read channel, closing it:", readErr)
			return
		}

		// Forward to all subscriptions
		if err := c.localTrack.WriteRTP(packet); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			logger.Println("Something went wrong in channel", c.id, "of", c.track.id)
		}
	}
}
