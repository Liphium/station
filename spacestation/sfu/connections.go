package sfu

import (
	"sync"

	"github.com/pion/webrtc/v4"
)

// Room id -> *sync.Map( Client id -> *sfu.Client )
var members *sync.Map

type Client struct {
	mutex      *sync.Mutex
	connection *webrtc.PeerConnection
}

// Register a new WebRTC connection for a specific client.
//
// Returns an offer from the server.
func NewClientConnection(room string, client string, offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {

	// Create a new peer connection for the user
	peer, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:" + defaultStunServer},
			},
		},
	})
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Set the description of the client
	if err := peer.SetRemoteDescription(offer); err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Let the gateway handle the rest of the connection
	manageConnection(room, client, peer)

	// Create an answer for the client
	answer, err := peer.CreateAnswer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	if err := peer.SetLocalDescription(answer); err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Wait for the ICE gathering to be completed
	// TODO: Consider implementing trickle ice instead of this garbage
	gatherComplete := webrtc.GatheringCompletePromise(peer)
	<-gatherComplete

	return *peer.LocalDescription(), nil
}
