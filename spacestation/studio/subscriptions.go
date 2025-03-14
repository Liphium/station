package studio

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

type Subscription struct {
	mutex   *sync.Mutex
	client  *Client                // read-only
	track   *Track                 // read-only
	peer    *webrtc.PeerConnection // read-only
	sender  *webrtc.RTPSender      // read-only
	channel string                 // The channel the client is subscribed to
}

// Delete the subscription
func (s *Subscription) Delete() {

	// Remove the track from the connection
	if err := s.peer.RemoveTrack(s.sender); err != nil {
		logger.Println("WARNING: Couldn't remove track from client:", err)
	}

	// Prevent concurrent modification of channel
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Delete the subscription from the channel in case it exists
	if obj, valid := s.track.channels.Load(s.channel); valid {
		c := obj.(*Channel)
		c.subscriptions.Delete(s.client.id)
	}

	// Delete the subscription from the client in case it exists
	s.client.subscriptions.Delete(s.track.id)
}
