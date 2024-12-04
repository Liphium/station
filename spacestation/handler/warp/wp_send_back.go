package warp_handlers

import (
	"slices"

	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: wp_send_back
func sendPacketBack(c *pipeshandler.Context, action struct {
	Warp       string `json:"w"` // The id of the Warp
	Target     string `json:"t"` // The target receiver of the packet
	Connection uint   `json:"c"` // The id of the connection this goes to (sometimes multiple ones need to be proxied)
	Packet     string `json:"p"` // The TCP packet that needs to be sent through Warp
}) pipes.Event {

	// Get the Warp related to the packet
	warp, err := caching.GetWarp(c.Client.Session, action.Warp)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, err)
	}

	// Make sure it's the hoster sending the event
	if warp.Hoster != c.Client.ID {
		return pipeshandler.ErrorResponse(c, localization.ErrorNoPermission, err)
	}

	// Lock the mutex to ensure there are no concurrent reads/writes
	warp.Mutex.Lock()

	// Make sure the target is in the warp
	if !slices.Contains(warp.Receivers, action.Target) {
		warp.Mutex.Unlock()
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequestContent, err)
	}

	// Unlock the mutex as it's not needed anymore (hoster shouldn't change?)
	warp.Mutex.Unlock()

	// Send the event to the hoster through the event channel
	if err := caching.SSNode.SendClient(action.Target, pipes.Event{
		Name: "wp_back",
		Data: map[string]interface{}{
			"c": action.Connection,
			"p": action.Packet,
		},
	}); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.SuccessResponse(c)
}
