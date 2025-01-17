package studio

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

type Subscription struct {
	mutex     *sync.Mutex
	client    string                      // Id of the client subscribing
	track     string                      // Id of the track subscribed to
	channel   string                      // The channel the client is subscribed to
	sendTrack *webrtc.TrackLocalStaticRTP // Track sent to the client
}

func (s *Subscription) Delete() {

}
