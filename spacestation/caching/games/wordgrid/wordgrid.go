package wordgrid

import (
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/caching/games"
	games_actions "github.com/Liphium/station/spacestation/handler/games"
	"github.com/Liphium/station/spacestation/util"
)

var (
	gameStatePick = 2
)

func LaunchWordGrid(session string) chan games.EventContext {
	channel := make(chan games.EventContext)
	go func() {
		for {
			event := <-channel
			if event.Name == "close" && event.Client == nil {
				break
			}

			if event.Client == nil {
				handleSysEvents(session, event)
				continue
			}

			util.Log.Println(event.Name)
		}
	}()
	return channel
}

func handleSysEvents(sessionId string, event games.EventContext) {
	switch event.Name {
	case "start":
		session, valid := caching.SetGameState(sessionId, 2)
		if !valid {
			return
		}

		games_actions.SendSessionUpdate(session)
	}
}
