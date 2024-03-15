package games_actions

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: game_init
func initGame(ctx pipeshandler.Context) {

	if ctx.ValidateForm("game") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	gameId := ctx.Data["game"].(string)
	conn, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	// Close current session
	currentSession, valid := caching.GetSession(conn.CurrentSession)
	if valid && len(currentSession.ClientIds) == 1 {
		caching.CloseSession(currentSession.Id)
		valid := sendSessionClose(ctx.Client.Session, currentSession.Id)
		if !valid {
			pipeshandler.ErrorResponse(ctx, "server.error")
			return
		}
	}

	session, valid := caching.OpenGameSession(ctx.Client.ID, conn.ClientID, ctx.Client.Session, gameId)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	caching.JoinSession(ctx.Client.ID, session.Id)

	// Send new session to all clients
	clients, valid := caching.GetAllConnections(ctx.Client.Session)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}
	adapters := make([]string, len(clients))
	i := 0
	for _, client := range clients {
		adapters[i] = client.Adapter
		i++
	}

	err := sendUpdateSession(adapters, session)
	if err != nil {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.NormalResponse(ctx, map[string]interface{}{
		"success": true,
		"session": session.Id,
		"min":     caching.GamesMap[session.Game].MinPlayers,
		"max":     caching.GamesMap[session.Game].MaxPlayers,
	})
}
