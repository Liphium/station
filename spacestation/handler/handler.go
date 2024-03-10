package handler

import (
	"github.com/Liphium/station/pipeshandler/wshandler"
	games_actions "github.com/Liphium/station/spacestation/handler/games"
	tabletop_handlers "github.com/Liphium/station/spacestation/handler/tabletop"
)

func Initialize() {
	wshandler.Initialize()
	games_actions.SetupActions()

	wshandler.Routes["set_data"] = setData
	wshandler.Routes["setup"] = setup
	wshandler.Routes["update"] = update

	tabletop_handlers.SetupHandler()
}
