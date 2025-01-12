package tabletop_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

type cursorMoveAction struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color float64 `json:"c"`
}

// Action: tc_move
func moveCursor(ctx *pipeshandler.Context, action cursorMoveAction) pipes.Event {

	// Notify other clients
	valid := caching.SendEventToMembers(ctx.Client.Session, pipes.Event{
		Name: "tc_moved",
		Data: map[string]interface{}{
			"c":   ctx.Client.ID,
			"x":   action.X,
			"y":   action.Y,
			"col": action.Color,
		},
	})
	if !valid {
		return pipeshandler.ErrorResponse(ctx, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(ctx)
}
