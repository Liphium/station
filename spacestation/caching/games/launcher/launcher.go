package launcher

import (
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/caching/games"
	"github.com/Liphium/station/spacestation/caching/games/wordgrid"
)

func InitGames() {
	caching.GamesMap["wordgrid"] = games.Game{
		Id:         "wordgrid",
		LaunchFunc: LaunchWorldGrid,
		MinPlayers: 1,
		MaxPlayers: 100,
	}
}

func LaunchWorldGrid(session string) chan games.EventContext {
	return wordgrid.LaunchWordGrid(session)
}
