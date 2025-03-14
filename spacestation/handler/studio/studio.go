package studio_handlers

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	studio_track_handlers "github.com/Liphium/station/spacestation/handler/studio/tracks"
)

func SetupHandlers() {
	pipeshandler.CreateHandlerFor(caching.SSInstance, "st_leave", leaveStudio)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "st_join", joinStudio)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "st_info", getStudioInfo)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "st_reneg", renegotiate)

	studio_track_handlers.SetupHandlers()
}
