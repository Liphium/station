package warp_handlers

import (
	"slices"

	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: wp_send_to
func sendPacket(c *pipeshandler.Context, action struct {
	Warp   string `json:"w"` // The id of the Warp
	Packet string `json:"p"` // The TCP packet that needs to be sent through Warp
}) pipes.Event {

	// Get the Warp related to the packet
	warp, err := caching.GetWarp(c.Client.Session, action.Warp)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, err)
	}

	// Lock the mutex to ensure there are no concurrent reads/writes
	warp.Mutex.Lock()

	// Check if the client is already in there
	if !slices.Contains(warp.Receivers, c.Client.ID) {

		// Add the client to the warp
		warp.Receivers = append(warp.Receivers, c.Client.ID)
	}

	// Unlock the mutex as it's not needed anymore (hoster shouldn't change?)
	warp.Mutex.Unlock()

	// Send the event to the hoster through the event channel
	if err := caching.SSNode.SendClient(warp.Hoster, pipes.Event{
		Name: "wp_to",
		Data: map[string]interface{}{
			"p": action.Packet,
		},
	}); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.SuccessResponse(c)
}
