package studio

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

type Track struct {
	id            string // Id of the track
	sender        string // Client id of the sender
	mutex         *sync.Mutex
	paused        bool
	simulcast     bool
	tracks        map[string]*webrtc.TrackRemote // Multiple tracks in case of simulcasting
	subscriptions []*Subscription
}

func (t *Track) IsPaused() bool {
	return t.paused
}

type Subscription struct {
}
