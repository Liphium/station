package games_actions

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: game_start
func startGame(ctx pipeshandler.Context) {

	if ctx.ValidateForm("session") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	sessionId := ctx.Data["session"].(string)
	session, valid := caching.GetSession(sessionId)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	if session.Creator != ctx.Client.ID {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	valid = caching.StartGameSession(sessionId)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "no.start")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}
