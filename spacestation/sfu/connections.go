package sfu

import (
	"errors"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
)

// Errors
var (
	clientNotFoundErr = errors.New("the specified client wasn't found")
)

// Room id -> *sync.Map( Client id -> *sfu.Client )
var roomMap *sync.Map

type Client struct {
	mutex         *sync.Mutex
	connection    *webrtc.PeerConnection
	lastKeepAlive time.Time
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

	// Remove the old connection in case there
	clientMap := getClientMap(room)
	if obj, valid := clientMap.Load(client); valid {
		oldConn := obj.(*Client)
		oldConn.mutex.Lock()
		oldConn.connection.Close()
		oldConn.mutex.Unlock()
	}

	// Add the new connection for the client
	clientMap.Store(client, &Client{
		mutex:      &sync.Mutex{},
		connection: peer,
	})

	return *peer.LocalDescription(), nil
}

// Update the keep alive
func updateKeepAlive(room string, client string) error {
	clientMap := getClientMap(room)
	obj, valid := clientMap.Load(client)
	if !valid {
		return clientNotFoundErr
	}
	cl := obj.(*Client)

	// Update the keep alive time to now
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.lastKeepAlive = time.Now()

	return nil
}

func Disconnect(room string, client string) error {
	clientMap := getClientMap(room)
	obj, valid := clientMap.Load(client)
	if !valid {
		return clientNotFoundErr
	}
	cl := obj.(*Client)

	// Disconnect the client
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.connection.Close()

	return nil
}

// Get the map of all client connections for a room
func getClientMap(room string) *sync.Map {
	var clientMap *sync.Map
	obj, valid := roomMap.Load(room)
	if !valid {
		clientMap = &sync.Map{}
		roomMap.Store(room, clientMap)
	} else {
		clientMap = obj.(*sync.Map)
	}
	return clientMap
}
