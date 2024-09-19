package handler

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	tabletop_handlers "github.com/Liphium/station/spacestation/handler/tabletop"
)

func Initialize() {
	pipeshandler.CreateHandlerFor(caching.SSInstance, "setup", setup)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "update", update)

	tabletop_handlers.SetupHandler()
}
