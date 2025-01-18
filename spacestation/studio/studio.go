package studio

import (
	"errors"
	"sync"

	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/pion/webrtc/v4"
)

// Errors
var (
	ErrClientNotFound = errors.New("the specified client wasn't found")
)

// Room id -> *Studio
var studioMap *sync.Map

type Studio struct {
	room    string
	clients *sync.Map // Client Id -> *Client
	tracks  *sync.Map // Track Id -> *Track
}

type Client struct {
	id              string                 // read-only
	connection      *webrtc.PeerConnection // read-only
	publishedTracks *sync.Map              // Track id (from client) -> *Track (read-only)
	subscriptions   *sync.Map              // Track id (server) -> *Subscription (read-only)
}

// Register a new WebRTC connection for a specific client.
//
// Returns an offer from the server.
func (s *Studio) NewClientConnection(client string, offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {

	// Create a new peer connection for the user
	peer, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:" + DefaultStunServer},
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
		id:              client,
		connection:      peer,
		publishedTracks: &sync.Map{},
		subscriptions:   &sync.Map{},
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

// Disconnect a client from studio
func (s *Studio) Disconnect(client string) error {
	obj, valid := s.clients.LoadAndDelete(client)
	if !valid {
		return ErrClientNotFound
	}
	cl := obj.(*Client)

	// Disconnect the client
	cl.connection.Close()

	return nil
}

// Send an event to everyone in the studio
func (s *Studio) SendEventToAll(event pipes.Event) error {

	// Get a list of all the adapters of the clients
	adapters := []string{}
	s.clients.Range(func(_, value any) bool {
		member := value.(*Client)
		adapters = append(adapters, member.id)
		return true
	})

	// Send the event through pipes
	if err := caching.SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(adapters),
		Local:   true,
		Event:   event,
	}); err != nil {
		logger.Println("error during event sending to studio members:", err)
		return err
	}

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
