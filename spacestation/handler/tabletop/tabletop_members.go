package tabletop_handlers

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: table_join
func joinTable(c *pipeshandler.Context, action struct {
	Color float64 `json:"color"`
}) pipes.Event {

	err := caching.JoinTable(c.Client.Session, c.Client.ID, action.Color)
	if err != nil {
		util.Log.Println("Couldn't join table of room", c.Client.Session, ":", err.Error())
		return pipeshandler.ErrorResponse(c, "server.error", err)
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

// Action: table_leave
func leaveTable(c *pipeshandler.Context, action interface{}) pipes.Event {
	err := caching.LeaveTable(c.Client.Session, c.Client.ID)
	if err != nil {
		util.Log.Println("Couldn't leave table of room", c.Client.Session, ":", err.Error())
		return pipeshandler.ErrorResponse(c, "server.error", err)
	}

	return pipeshandler.SuccessResponse(c)
}
