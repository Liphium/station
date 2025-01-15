package sfu

import (
	"sync"

	"github.com/Liphium/station/spacestation/util"
	"github.com/pion/ice/v4"
	"github.com/pion/webrtc/v4"
)

// The base webrtc API used for all connections
var api *webrtc.API

func Start(port int) {

	// Create a new setting engine
	engine := webrtc.SettingEngine{}

	// Set the port
	mux, err := ice.NewMultiUDPMuxFromPort(port)
	if err != nil {
		util.Log.Fatal("Couldn't create port multiplexer for the SFU:", err)
	}
	engine.SetICEUDPMux(mux)

	// Create the api using the settings engine
	api = webrtc.NewAPI(webrtc.WithSettingEngine(engine))
}

// Room id -> *sync.Map( Client id -> Member info )
var members *sync.Map

type MemberInfo struct {
	// Whether or not the client is actually connected to the WebRTC room
	Connected bool

	// For preventing clashes between multiple goroutines
	Mutex *sync.Mutex

	// All the tracks the client is publishing
	Tracks []webrtc.TrackRemote
}
