package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

// Action: tc_move
func moveCursor(ctx pipeshandler.Context) {

	if ctx.ValidateForm("x", "y") {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get the connection (for getting the client id)
	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	// Get all the data needed
	x := ctx.Data["x"].(float64)
	y := ctx.Data["y"].(float64)
	member, valid := caching.GetMemberData(ctx.Client.Session, ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	// Notify other clients
	valid = SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tc_moved",
		Data: map[string]interface{}{
			"c":   connection.ClientID,
			"x":   x,
			"y":   y,
			"col": member.Color,
		},
	})
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}
