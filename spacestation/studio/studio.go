package studio

import (
	"errors"
	"sync"

	"github.com/pion/webrtc/v4"
)

// Errors
var (
	errClientNotFound = errors.New("the specified client wasn't found")
)

// Room id -> *Studio
var studioMap *sync.Map

type Studio struct {
	room    string
	clients *sync.Map // Client Id -> *Client
	tracks  *sync.Map // Track Id -> *Track
}

type Client struct {
	mutex           *sync.Mutex
	connection      *webrtc.PeerConnection
	publishedTracks []*Track
	subscriptions   []*Subscription
}

// Register a new WebRTC connection for a specific client.
//
// Returns an offer from the server.
func (s *Studio) NewClientConnection(room string, client string, offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {

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

	// Disconnect the old connection in case there
	if _, valid := s.clients.Load(client); valid {
		if err := s.Disconnect(client); err != nil {
			return webrtc.SessionDescription{}, err
		}
	}

	// Add the new connection for the client
	c := &Client{
		mutex:      &sync.Mutex{},
		connection: peer,
	}
	s.clients.Store(client, c)

	// Let the gateway handle the rest of the connection
	s.startGateway(c, peer)

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

func (s *Studio) Disconnect(client string) error {
	obj, valid := s.clients.Load(client)
	if !valid {
		return errClientNotFound
	}
	cl := obj.(*Client)

	// Disconnect the client
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.connection.Close()

	return nil
}

// Get the studio for a specific room
func GetStudio(room string) *Studio {
	var clientMap *Studio
	obj, valid := studioMap.Load(room)
	if !valid {
		clientMap = &Studio{
			room:    room,
			clients: &sync.Map{},
			tracks:  &sync.Map{},
		}
		studioMap.Store(room, clientMap)
	} else {
		clientMap = obj.(*Studio)
	}
	return clientMap
}
