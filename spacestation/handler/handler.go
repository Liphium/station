package handler

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	message_handlers "github.com/Liphium/station/spacestation/handler/messages"
	tabletop_handlers "github.com/Liphium/station/spacestation/handler/tabletop"
	warp_handlers "github.com/Liphium/station/spacestation/handler/warp"
)

func Initialize() {
	pipeshandler.CreateHandlerFor(caching.SSInstance, "setup", setup)

	message_handlers.SetupHandler()
	tabletop_handlers.SetupHandler()
	warp_handlers.SetupHandlers()
}
