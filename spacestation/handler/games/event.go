package games_actions

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/caching/games"
)

// Action: game_event
func gameEvent(ctx pipeshandler.Context) {

	if ctx.ValidateForm("session", "name", "data") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	sessionId := ctx.Data["session"].(string)
	conn, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	if conn.CurrentSession != sessionId {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	valid = caching.ForwardGameEvent(sessionId, games.EventContext{
		Client: ctx.Client,
		Name:   ctx.Data["name"].(string),
		Data:   ctx.Data["data"],
	})
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}
