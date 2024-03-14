package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler/wshandler"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandler() {

	// Table member management
	wshandler.RegisterHandler(caching.Node, "table_join", joinTable)
	wshandler.RegisterHandler(caching.Node, "table_leave", leaveTable)

	// Table object management
	wshandler.RegisterHandler(caching.Node, "tobj_create", createObject)
	wshandler.RegisterHandler(caching.Node, "tobj_delete", deleteObject)
	wshandler.RegisterHandler(caching.Node, "tobj_select", selectObject)
	wshandler.RegisterHandler(caching.Node, "tobj_modify", modifyObject)
	wshandler.RegisterHandler(caching.Node, "tobj_move", moveObject)
	wshandler.RegisterHandler(caching.Node, "tobj_rotate", rotateObject)

	// Table cursor sending
	wshandler.RegisterHandler(caching.Node, "tc_move", moveCursor)
}

// Send an event to all table members
func SendEventToMembers(room string, event pipes.Event) bool {
	valid, members := caching.TableMembers(room)
	if !valid {
		return false
	}

	return caching.Node.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(members),
		Local:   true,
		Event:   event,
	}) == nil
}
