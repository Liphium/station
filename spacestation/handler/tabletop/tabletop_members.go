package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: table_enable
func enableTable(c *pipeshandler.Context, _ interface{}) pipes.Event {

	// Enable the member
	if err := caching.ChangeTableMemberState(c.Client.Session, c.Client.ID, true); err != nil {
		return pipeshandler.ErrorResponse(c, err.Error(), err)
	}

	// Start a goroutine to stream over all the changes
	go func(c pipeshandler.Context) {
		// Send all the objects
		objects, err := caching.TableObjects(c.Client.Session)
		if err != nil {
			util.Log.Println("Couldn't get objects of room", c.Client.Session, ":", err.Error())
			return
		}

		err = caching.SSNode.SendClient(c.Client.ID, pipes.Event{
			Name: "table_obj",
			Data: map[string]interface{}{
				"obj": objects,
			},
		})
		if err != nil {
			util.Log.Println("Couldn't send objects of room through event", c.Client.Session, ":", err.Error())
		}
	}(*c)

	return pipeshandler.SuccessResponse(c)
}

// Action: table_disable
func disableTable(c *pipeshandler.Context, action interface{}) pipes.Event {
	err := caching.ChangeTableMemberState(c.Client.Session, c.Client.ID, false)
	if err != nil {
		util.Log.Println("Couldn't disable table of room", c.Client.Session, "for", c.Client.ID, ":", err.Error())
		return pipeshandler.ErrorResponse(c, "server.error", err)
	}

	return pipeshandler.SuccessResponse(c)
}
