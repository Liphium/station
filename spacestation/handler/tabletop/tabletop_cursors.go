package tabletop_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

type cursorMoveAction struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Action: tc_move
func moveCursor(ctx *pipeshandler.Context, action cursorMoveAction) pipes.Event {

	// Get the connection (for getting the client id)
	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(ctx, localization.ErrorInvalidRequest, nil)
	}

	// Get all the data needed
	member, valid := caching.GetMemberData(ctx.Client.Session, ctx.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(ctx, localization.ErrorServer, nil)
	}

	// Notify other clients
	valid = SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tc_moved",
		Data: map[string]interface{}{
			"c":   connection.ClientID,
			"x":   action.X,
			"y":   action.Y,
			"col": member.Color,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(ctx, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(ctx)
}
