package studio

import (
	"strings"
	"sync"

	"github.com/Liphium/station/spacestation/util"
	"github.com/pion/webrtc/v4"
)

type Client struct {
	id              string                 // read-only
	studio          *Studio                // read-only
	connection      *webrtc.PeerConnection // read-only
	publishedTracks *sync.Map              // Track id (from client) -> *Track (read-only)
	subscriptions   *sync.Map              // Track id (server) -> *Subscription (read-only)
}

// Get a client's subscription to a specific track
func (c *Client) GetSubscription(track string) (*Subscription, bool) {
	obj, valid := c.subscriptions.Load(track)
	if !valid {
		return nil, false
	}
	return obj.(*Subscription), true
}

// Initialize the handlers for the client's connection
func (c *Client) initializeConnection(peer *webrtc.PeerConnection) error {

	// Listen for any data channels
	peer.OnDataChannel(func(dc *webrtc.DataChannel) {
		logger.Println("new data channel", dc.Label())
	})

	// Check if negotiation is needed
	peer.OnNegotiationNeeded(func() {
		logger.Println("renegotiation is needed")
	})

	// Listen for new tracks
	peer.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {

		logger.Println("received new track", tr.ID(), tr.RID())

		// Parse the track rid for the channel
		channel := tr.RID()
		if channel == "" {
			channel = "f"
		}

		// Make sure the channel is valid
		if !strings.Contains("fhq", channel) {
			logger.Println(c.id, "disconnected due to wrong channel:", channel)
			c.studio.Disconnect(c.id)
			return
		}

		// Check if the track has already been published with this id
		var track *Track
		if obj, valid := c.publishedTracks.Load(tr.ID()); valid {

			// Add the channel
			track = obj.(*Track)
			track.AddChannel(channel, tr)
		} else {

			// Generate a new id for the track
			id := util.GenerateToken(12)
			_, valid := c.studio.tracks.Load(id)
			for valid {
				id = util.GenerateToken(12)
				_, valid = c.studio.tracks.Load(id)
			}

			// Register the track
			track = &Track{
				studio:        c.studio,
				id:            id,
				mutex:         &sync.Mutex{},
				sender:        c.id,
				senderTrack:   tr.ID(),
				paused:        false,
				simulcast:     false,
				channelAmount: 0,
				channels:      &sync.Map{},
			}
			c.studio.tracks.Store(id, track)

			// Add the channel in all places
			track.AddChannel(channel, tr)
			c.publishedTracks.Store(tr.ID(), track)
		}
	})

	// Disconnect the client when the connection closes
	peer.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateClosed {
			c.studio.Disconnect(c.id)
		}
	})

	// Send them all the tracks currently available
	go func() {
		c.studio.tracks.Range(func(key, value any) bool {
			t := value.(*Track)
			t.SendTrackUpdate(c.id, true)
			return true
		})
	}()

	return nil
}

// Handle the removing of a track the client is sending to studio
func (c *Client) handleRemoveTrack(t *Track) {

	// Delete the track from the published tracks
	_, valid := c.publishedTracks.LoadAndDelete(t.id)
	if !valid {
		return
	}

	// End the receiver currently receiving the track
	for _, r := range c.connection.GetReceivers() {
		if r.Track().ID() == t.senderTrack {
			if err := r.Stop(); err != nil {
				logger.Println("warning: couldn't stop receiver of ended track")
			}
		}
	}
}
