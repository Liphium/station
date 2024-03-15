package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandler() {

	// Table member management
	caching.Instance.RegisterHandler("table_join", joinTable)
	caching.Instance.RegisterHandler("table_leave", leaveTable)

	// Table object management
	caching.Instance.RegisterHandler("tobj_create", createObject)
	caching.Instance.RegisterHandler("tobj_delete", deleteObject)
	caching.Instance.RegisterHandler("tobj_select", selectObject)
	caching.Instance.RegisterHandler("tobj_modify", modifyObject)
	caching.Instance.RegisterHandler("tobj_move", moveObject)
	caching.Instance.RegisterHandler("tobj_rotate", rotateObject)

	// Table cursor sending
	caching.Instance.RegisterHandler("tc_move", moveCursor)
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
