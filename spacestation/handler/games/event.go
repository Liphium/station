package games_actions

import (
	"github.com/Liphium/station/pipeshandler/wshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/caching/games"
)

// Action: game_event
func gameEvent(message wshandler.Message) {

	if message.ValidateForm("session", "name", "data") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	sessionId := message.Data["session"].(string)
	conn, valid := caching.GetConnection(message.Client.ID)
	if !valid {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	if conn.CurrentSession != sessionId {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	valid = caching.ForwardGameEvent(sessionId, games.EventContext{
		Client: message.Client,
		Name:   message.Data["name"].(string),
		Data:   message.Data["data"],
	})
	if !valid {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	wshandler.SuccessResponse(message)
}
