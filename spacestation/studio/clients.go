package studio

import (
	"slices"
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

	// Create a new data channel for pipes (this is gonna be used in the future)
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

		// Check if the channel is valid (RID specifies channel)
		if !slices.Contains(acceptedChannels, tr.RID()) {
			logger.Println(c.id+"disconnected due to wrong channel (", tr.RID(), ")")
			c.studio.Disconnect(c.id)
			return
		}

		// Check if the track has already been published with this id
		var track *Track
		if obj, valid := c.publishedTracks.Load(tr.ID()); valid {

			// Add the channel
			track = obj.(*Track)
			track.AddChannel(tr)
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
				studio:      c.studio,
				id:          id,
				sender:      c.id,
				senderTrack: tr.ID(),
				mutex:       &sync.Mutex{},
				paused:      false,
				simulcast:   false,
				channels:    &sync.Map{},
			}
			c.studio.tracks.Store(id, track)

			// Add the channel in all places
			track.AddChannel(tr)
			c.publishedTracks.Store(tr.ID(), track)
		}
	})

	return nil
}