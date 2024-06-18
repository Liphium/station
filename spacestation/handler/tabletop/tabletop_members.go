package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: table_join
func joinTable(ctx pipeshandler.Context) {

	err := caching.JoinTable(ctx.Client.Session, ctx.Client.ID)
	if err != nil {
		util.Log.Println("Couldn't join table of room", ctx.Client.Session, ":", err.Error())
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)

	// Send all the objects
	objects, err := caching.TableObjects(ctx.Client.Session)
	if err != nil {
		util.Log.Println("Couldn't get objects of room", ctx.Client.Session, ":", err.Error())
		return
	}

	err = caching.SSNode.SendClient(ctx.Client.ID, pipes.Event{
		Name: "table_obj",
		Data: map[string]interface{}{
			"obj": objects,
		},
	})
	if err != nil {
		util.Log.Println("Couldn't send objects of room through event", ctx.Client.Session, ":", err.Error())
	}
}

// Action: table_leave
func leaveTable(ctx pipeshandler.Context) {
	err := caching.LeaveTable(ctx.Client.Session, ctx.Client.ID)
	if err != nil {
		util.Log.Println("Couldn't leave table of room", ctx.Client.Session, ":", err.Error())
		pipeshandler.ErrorResponse(ctx, "server.error")
		return
	}

	pipeshandler.SuccessResponse(ctx)
}
