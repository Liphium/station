package studio

import (
	"slices"
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/caching"
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
	id          string  // Id of the track (read-only)
	studio      *Studio // The studio the track belongs to
	sender      string  // Client id of the sender
	senderTrack string  // How the client sending refers to this track
	mutex       *sync.Mutex
	paused      bool
	simulcast   bool
	channels    *sync.Map // Channel id -> *Channel
}

// Add a new channel for a track
func (t *Track) AddChannel(tr *webrtc.TrackRemote) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Add as a new channel
	channel := &Channel{
		track:       t,
		id:          tr.RID(),
		remoteTrack: tr,
	}
	if !t.simulcast {

		// If the channel has a different id than any previous channel, turn on simulcast
		_, valid := t.channels.Load(tr.RID())
		t.simulcast = !valid
	}
	t.channels.Store(tr.RID(), channel)

	// Start the sender
	go channel.startSender()

	// Send update notifying everyone of the change
	if err := t.SendTrackUpdateToAll(false); err != nil {
		logger.Println("WARNING: Couldn't send track update:", err)
	}
}

// Send an event that updates the track to one client
//
// If mutex is true, the track's mutex will be locked.
func (t *Track) SendTrackUpdate(client string, mutex bool) error {
	return caching.SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.P2PChannel(client, caching.SSNode.ID),
		Event:   t.TrackUpdateEvent(mutex),
	})
}

// Send an event that updates the track on all clients in studio.
//
// If mutex is true, the track's mutex will be locked.
func (t *Track) SendTrackUpdateToAll(mutex bool) error {

	// Send the updated track to everyone
	return t.studio.SendEventToAll(t.TrackUpdateEvent(mutex))
}

// Get the event required for updating the track.
//
// If mutex is true, the track's mutex will be locked.
func (t *Track) TrackUpdateEvent(mutex bool) pipes.Event {

	// Get all the subscribers of the track
	var channels []string
	var subscribers []string
	t.channels.Range(func(key, value any) bool {
		c := value.(*Channel)

		// Add all subscribers from the channel
		c.subscriptions.Range(func(key, value any) bool {
			clientId := key.(string)
			if !slices.Contains(subscribers, clientId) {
				subscribers = append(subscribers, clientId)
			}
			return true
		})

		// Add the channel itself
		channels = append(channels, c.id)

		return true
	})

	// Make sure there are no concurrent modifications
	if mutex {
		t.mutex.Lock()
		defer t.mutex.Unlock()
	}

	// Send the updated track to everyone
	return pipes.Event{
		Name: "st_tr_update",
		Data: map[string]interface{}{
			"track":    t.id,
			"sender":   t.sender,
			"paused":   t.paused,
			"channels": channels,
			"subs":     subscribers,
		},
	}
}

// Create a new subscription for a specific channel
func (t *Track) NewSubscription(c *Client, channel string) {

	// Get the channel

	// Create a new local track
	track, err := webrtc.NewTrackLocalStaticRTP()

	// Get the connection for the client

	sub := &Subscription{
		mutex:   &sync.Mutex{},
		client:  c.id,
		track:   t.id,
		channel: channel,
	}

}

func (t *Track) IsPaused() bool {
	return t.paused
}
