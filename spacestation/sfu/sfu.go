package sfu

import (
	"log"
	"os"

	"github.com/Liphium/station/spacestation/util"
	"github.com/pion/ice/v4"
	"github.com/pion/webrtc/v4"
)

// The base webrtc API used for all connections
var api *webrtc.API

// Custom logger for everything going on in the SFU
var logger *log.Logger = log.New(os.Stdout, "space-sfu ", log.Flags())

// Configuration
var Enabled = false // Changed later in the setup
var defaultStunServer = "stun.l.google.com:19302"

// TODO: Add turn server support

func Start(port int) {

	if os.Getenv("SS_SFU_ENABLE") != "" {
		if os.Getenv("SS_SFU_ENABLE") == "false" {
			logger.Println("SFU disabled, as request by the SS_SFU_ENABLE environment variable.")
			return
		}

		logger.Println("Starting SFU..")
		Enabled = true
	} else {
		logger.Println("Spaces SFU not enabled. Voice and video for Spaces has been disabled.")
		return
	}

	// Get all the environment variables
	if os.Getenv("SS_STUN") != "" {
		defaultStunServer = os.Getenv("SS_STUN")
	} else {
		logger.Println("WARNING: No STUN server provided, using Google's one instead. Read more at https://docs.liphium.com/setup/config-setup.")
	}

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

	logger.Println("SFU started. Voice and video for Spaces has been enabled.")
}
