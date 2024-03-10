package games_actions

import (
	"github.com/Liphium/station/pipeshandler/wshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: game_start
func startGame(message wshandler.Message) {

	if message.ValidateForm("session") {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	sessionId := message.Data["session"].(string)
	session, valid := caching.GetSession(sessionId)
	if !valid {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	if session.Creator != message.Client.ID {
		wshandler.ErrorResponse(message, "invalid")
		return
	}

	valid = caching.StartGameSession(sessionId)
	if !valid {
		wshandler.ErrorResponse(message, "no.start")
		return
	}

	wshandler.SuccessResponse(message)
}
