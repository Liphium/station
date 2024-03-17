package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandler() {

	// Table member management
	caching.SSInstance.RegisterHandler("table_join", joinTable)
	caching.SSInstance.RegisterHandler("table_leave", leaveTable)

	// Table object management
	caching.SSInstance.RegisterHandler("tobj_create", createObject)
	caching.SSInstance.RegisterHandler("tobj_delete", deleteObject)
	caching.SSInstance.RegisterHandler("tobj_select", selectObject)
	caching.SSInstance.RegisterHandler("tobj_modify", modifyObject)
	caching.SSInstance.RegisterHandler("tobj_move", moveObject)
	caching.SSInstance.RegisterHandler("tobj_rotate", rotateObject)

	// Table cursor sending
	caching.SSInstance.RegisterHandler("tc_move", moveCursor)
}

// Send an event to all table members
func SendEventToMembers(room string, event pipes.Event) bool {
	valid, members := caching.TableMembers(room)
	if !valid {
		return false
	}

	return caching.SSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(members),
		Local:   true,
		Event:   event,
	}) == nil
}
