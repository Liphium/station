package handler

import (
	"github.com/Liphium/station/chatserver/caching"
	games_actions "github.com/Liphium/station/spacestation/handler/games"
	tabletop_handlers "github.com/Liphium/station/spacestation/handler/tabletop"
)

func Initialize() {
	games_actions.SetupActions()

	caching.Instance.RegisterHandler("set_data", setData)
	caching.Instance.RegisterHandler("setup", setup)
	caching.Instance.RegisterHandler("update", update)

	tabletop_handlers.SetupHandler()
}
