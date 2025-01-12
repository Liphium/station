package tabletop_handlers

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandler() {

	// Table member management
	pipeshandler.CreateHandlerFor(caching.SSInstance, "table_enable", enableTable)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "table_disable", disableTable)

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
