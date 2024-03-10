package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/send"
	"github.com/Liphium/station/pipeshandler/wshandler"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandler() {

	// Table member management
	wshandler.Routes["table_join"] = joinTable
	wshandler.Routes["table_leave"] = leaveTable

	// Table object management
	wshandler.Routes["tobj_create"] = createObject
	wshandler.Routes["tobj_delete"] = deleteObject
	wshandler.Routes["tobj_select"] = selectObject
	wshandler.Routes["tobj_modify"] = modifyObject
	wshandler.Routes["tobj_move"] = moveObject
	wshandler.Routes["tobj_rotate"] = rotateObject

	// Table cursor sending
	wshandler.Routes["tc_move"] = moveCursor
}

// Send an event to all table members
func SendEventToMembers(room string, event pipes.Event) bool {
	valid, members := caching.TableMembers(room)
	if !valid {
		return false
	}

	return send.Pipe(send.ProtocolWS, pipes.Message{
		Channel: pipes.BroadcastChannel(members),
		Local:   true,
		Event:   event,
	}) == nil
}
