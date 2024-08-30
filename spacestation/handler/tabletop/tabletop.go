package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

func SetupHandler() {

	// Table member management
	pipeshandler.CreateHandlerFor(caching.SSInstance, "table_join", joinTable)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "table_leave", leaveTable)

	// Table object management
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_create", createObject)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_delete", deleteObject)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_select", selectObject)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_unselect", unselectObject)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_modify", modifyObject)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_move", moveObject)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_rotate", rotateObject)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tobj_mqueue", queueModificationToObject)

	// Table cursor sending
	pipeshandler.CreateHandlerFor(caching.SSInstance, "tc_move", moveCursor)
}

// Send an event to all table members
func SendEventToMembers(room string, event pipes.Event) bool {
	valid := caching.RangeOverTableMembers(room, func(tm *caching.TableMember) bool {
		if err := caching.SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
			Channel: pipes.BroadcastChannel([]string{tm.Client}),
			Local:   true,
			Event:   event,
		}); err != nil {
			util.Log.Println("error during event sending to tabletop members:", err)
		}
		return true
	})
	return valid
}
