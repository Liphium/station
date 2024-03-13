package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler/wshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: table_join
func joinTable(message wshandler.Message) {

	err := caching.JoinTable(message.Client.Session, message.Client.ID)
	if err != nil {
		util.Log.Println("Couldn't join table of room", message.Client.Session, ":", err.Error())
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	wshandler.SuccessResponse(message)

	// Send all the objects
	objects, err := caching.TableObjects(message.Client.Session)
	if err != nil {
		util.Log.Println("Couldn't get objects of room", message.Client.Session, ":", err.Error())
		return
	}

	err = message.Client.SendEvent(pipes.Event{
		Name: "table_obj",
		Data: map[string]interface{}{
			"obj": objects,
		},
	})
	if err != nil {
		util.Log.Println("Couldn't send objects of room through event", message.Client.Session, ":", err.Error())
	}
}

// Action: table_leave
func leaveTable(message wshandler.Message) {
	err := caching.LeaveTable(message.Client.Session, message.Client.ID)
	if err != nil {
		util.Log.Println("Couldn't leave table of room", message.Client.Session, ":", err.Error())
		wshandler.ErrorResponse(message, "server.error")
		return
	}

	wshandler.SuccessResponse(message)
}
