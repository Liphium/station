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
var StunServer = ""
var TurnServer = ""
var TurnUsername = ""
var TurnPassword = ""
var Port int = 0

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
	if os.Getenv("SS_SFU_PORT") != "" {
		var err error
		Port, err = strconv.Atoi(os.Getenv("SS_SFU_PORT"))
		if err != nil {
			logger.Fatal("Invalid port number in SS_PORT environment variable")
		}
	} else {
		Port = 5000
	}
	logger.Println("Starting on port", Port, "..")

	// Get all the environment variables
	if os.Getenv("SS_STUN") != "" {
		StunServer = os.Getenv("SS_STUN")
	} else {
		logger.Println("ERROR: No STUN server provided, can't start Liphium like this. Read more at https://docs.liphium.com/setup/config-setup.")
		Enabled = false
		return
	}
	if os.Getenv("SS_TURN") != "" {
		TurnServer = os.Getenv("SS_TURN")
		TurnUsername = os.Getenv("SS_TURN_USERNAME")
		TurnPassword = os.Getenv("SS_TURN_PASSWORD")
	} else {
		logger.Println("WARNING: No TURN server provided, this can cause connectivity issues for some people. Read more at https://docs.liphium.com/setup/config-setup.")
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
