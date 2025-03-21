package handler

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	message_handlers "github.com/Liphium/station/spacestation/handler/messages"
	studio_handlers "github.com/Liphium/station/spacestation/handler/studio"
	tabletop_handlers "github.com/Liphium/station/spacestation/handler/tabletop"
	warp_handlers "github.com/Liphium/station/spacestation/handler/warp"
)

func Initialize() {
	pipeshandler.CreateHandlerFor(caching.SSInstance, "setup", setup)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "set_audio_state", setAudioState)

	message_handlers.SetupHandler()
	tabletop_handlers.SetupHandler()
	warp_handlers.SetupHandlers()
	studio_handlers.SetupHandlers()
}
