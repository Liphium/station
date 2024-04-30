package handler

import (
	"github.com/Liphium/station/spacestation/caching"
	games_actions "github.com/Liphium/station/spacestation/handler/games"
	tabletop_handlers "github.com/Liphium/station/spacestation/handler/tabletop"
)

func Initialize() {
	games_actions.SetupActions()

	caching.SSInstance.RegisterHandler("setup", setup)
	caching.SSInstance.RegisterHandler("update", update)

	tabletop_handlers.SetupHandler()
}
