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

	connection, valid := caching.GetConnection(ctx.Client.ID)
	if !valid {
		pipeshandler.ErrorResponse(ctx, "invalid")
		return
	}

	x := ctx.Data["x"].(float64)
	y := ctx.Data["y"].(float64)

	// Notify other clients
	valid = SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tc_moved",
		Data: map[string]interface{}{
			"c": connection.ClientID,
			"x": x,
			"y": y,
		},
	})
	if !valid {
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}
