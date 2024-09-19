package tabletop_handlers

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: table_enable
func enableTable(c *pipeshandler.Context, _ interface{}) pipes.Event {

	// Enable the member
	if msg := caching.ChangeTableMemberState(c.Client.Session, c.Client.ID, true); msg != nil {
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Start a goroutine to stream over all the changes
	go func(c pipeshandler.Context) {
		// Send all the objects
		objects, msg := caching.TableObjects(c.Client.Session)
		if msg != nil {
			util.Log.Println("Couldn't get objects of room", c.Client.Session, ":", msg[localization.DefaultLocale])
			return
		}

		err := caching.SSNode.SendClient(c.Client.ID, pipes.Event{
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
	msg := caching.ChangeTableMemberState(c.Client.Session, c.Client.ID, false)
	if msg != nil {
		util.Log.Println("Couldn't disable table of room", c.Client.Session, "for", c.Client.ID, ":", msg[localization.DefaultLocale])
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	return pipeshandler.SuccessResponse(c)
}
