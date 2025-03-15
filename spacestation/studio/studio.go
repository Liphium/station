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
var studioMap *sync.Map = &sync.Map{}

type Studio struct {
	room    string
	clients *sync.Map // Client Id -> *Client
	tracks  *sync.Map // Track Id -> *Track
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

	/*
		// Add the required transceivers
		if _, err := peer.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
			return webrtc.SessionDescription{}, err
		}
	*/

	// Add the new connection for the client
	c := &Client{
		id:              client,
		studio:          s,
		connection:      peer,
		publishedTracks: &sync.Map{},
		subscriptions:   &sync.Map{},
	}
	s.clients.Store(client, c)

	// Let the gateway handle the rest of the connection
	if err := c.initializeConnection(peer); err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Create an answer for the client
	answer, err := peer.CreateAnswer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Set the local description
	if err := peer.SetLocalDescription(answer); err != nil {
		return webrtc.SessionDescription{}, err
	}

	return *peer.LocalDescription(), nil
}

// Renegotiate with a client.
//
// Returns the new answer for the client.
func (s *Studio) HandleClientRenegotiation(client string, offer webrtc.SessionDescription) (webrtc.SessionDescription, error) {

	// Get the client
	c, valid := s.GetClient(client)
	if !valid {
		return webrtc.SessionDescription{}, ErrClientNotFound
	}

	// Set the remote description
	if err := c.connection.SetRemoteDescription(offer); err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Create an answer for the client
	answer, err := c.connection.CreateAnswer(nil)
	if err != nil {
		return webrtc.SessionDescription{}, err
	}

	// Set the local description
	if err := c.connection.SetLocalDescription(answer); err != nil {
		return webrtc.SessionDescription{}, err
	}

	return *c.connection.LocalDescription(), nil
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
		client := value.(*Client)
		adapters = append(adapters, client.id)
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

// Forward a lightwire packet to all clients using it
func (s *Studio) ForwardLightwirePacket(packet []byte) {

	// Send it to all clients with Lightwire
	s.clients.Range(func(_, value any) bool {
		client := value.(*Client)
		if client.lightwire != nil {
			client.lightwire.SendPacket(packet)
		}
		return true
	})
}

// Get a track in the studio
func (s *Studio) GetTrack(track string) (*Track, bool) {
	obj, valid := s.tracks.Load(track)
	if !valid {
		return nil, false
	}
	return obj.(*Track), true
}

// Get a client in the studio
func (s *Studio) GetClient(client string) (*Client, bool) {
	obj, valid := s.clients.Load(client)
	if !valid {
		return nil, false
	}
	return obj.(*Client), true
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
