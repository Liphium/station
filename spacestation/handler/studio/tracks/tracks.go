package studio_track_handlers

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandlers() {
	pipeshandler.CreateHandlerFor(caching.SSInstance, "st_tr_subscribe", subscribeToTrack)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "st_tr_unsubscribe", unsubscribeToTrack)
}
