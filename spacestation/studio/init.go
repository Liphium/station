package studio

import (
	"log"
	"os"
	"strconv"

	"github.com/Liphium/station/spacestation/util"
	"github.com/pion/ice/v4"
	"github.com/pion/webrtc/v4"
)

// The base webrtc API used for all connections
var api *webrtc.API

// Custom logger for everything going on in the SFU
var logger *log.Logger = log.New(os.Stdout, "space-studio ", log.Flags())

// Configuration
var Enabled = false // Changed later in the setup
var DefaultStunServer = "stun.l.google.com:19302"
var Port int = 0

// TODO: Add turn server support

func Start() {

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

	// Set the port if available
	if os.Getenv("SS_PORT") != "" {
		var err error
		Port, err = strconv.Atoi(os.Getenv("SS_PORT"))
		if err != nil {
			logger.Fatal("Invalid port number in SS_PORT environment variable")
		}
	} else {
		Port = 5000
	}
	logger.Println("Starting on port", Port, "..")

	// Get all the environment variables
	if os.Getenv("SS_STUN") != "" {
		DefaultStunServer = os.Getenv("SS_STUN")
	} else {
		logger.Println("WARNING: No STUN server provided, using Google's one instead. Read more at https://docs.liphium.com/setup/config-setup.")
	}

	// Create a new setting engine
	engine := webrtc.SettingEngine{}

	// Set the port
	mux, err := ice.NewMultiUDPMuxFromPort(Port)
	if err != nil {
		util.Log.Fatal("Couldn't create port multiplexer for the SFU:", err)
	}
	engine.SetICEUDPMux(mux)

	// Create the api using the settings engine
	api = webrtc.NewAPI(webrtc.WithSettingEngine(engine))

	logger.Println("SFU started. Voice and video for Spaces has been enabled.")
}
