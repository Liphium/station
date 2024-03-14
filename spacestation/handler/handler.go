package handler

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipeshandler/wshandler"
	games_actions "github.com/Liphium/station/spacestation/handler/games"
	tabletop_handlers "github.com/Liphium/station/spacestation/handler/tabletop"
)

func Initialize() {
	wshandler.Initialize()
	games_actions.SetupActions()

	wshandler.RegisterHandler(caching.Node, "set_data", setData)
	wshandler.RegisterHandler(caching.Node, "setup", setup)
	wshandler.RegisterHandler(caching.Node, "update", update)

	tabletop_handlers.SetupHandler()
}
