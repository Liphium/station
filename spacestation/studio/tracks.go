package studio

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

// Accepted track channels
const (
	channelLow     = "l"
	channelMedium  = "m"
	channelHigh    = "h"
	channelDefault = "d"
)

// A list for filtering
var acceptedChannels = []string{channelLow, channelMedium, channelHigh, channelDefault}

type Track struct {
	id            string // Id of the track
	sender        string // Client id of the sender
	senderTrack   string // How the client sending refers to this track
	mutex         *sync.Mutex
	paused        bool
	simulcast     bool
	channels      map[string]*webrtc.TrackRemote // Channels (to allow things like simulcasting)
	subscriptions *sync.Map                      // Channel -> []*Subscription
}

// Add a new channel for a track
func (t *Track) AddChannel(channel string, tr *webrtc.TrackRemote) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Add as a new channel
	if !t.simulcast {
		t.simulcast = t.channels[tr.RID()] == nil
	}
	t.channels[tr.RID()] = tr

	// Start the sender
	go t.startSenderForChannel(channel, tr)
}

// Handles sending packets for a channel to all subscribers
func (t *Track) startSenderForChannel(channel string, tr *webrtc.TrackRemote) {
	for {
		// Read RTP packets being sent on the channel
		packet, _, readErr := tr.ReadRTP()
		if readErr != nil {
			panic(readErr) // Don't know what to do here yet
		}

		// Get all of the subscriptions for the current channel
		obj, valid := t.subscriptions.Load(channel)
		if !valid {
			continue
		}
		subs := obj.([]*Subscription)

		// Forward to all subscriptions
		for _, sub := range subs {
			sub.mutex.Lock()
			if err := sub.sendTrack.WriteRTP(packet); err != nil {
				logger.Println("Deleting subscription from", sub.client, "to", t.id, "(", t.senderTrack, ")", ":", err)
				sub.mutex.Unlock()
				sub.Delete()
				continue
			}
			sub.mutex.Unlock()
		}
	}
}

func (t *Track) IsPaused() bool {
	return t.paused
}
