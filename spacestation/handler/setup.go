package handler

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/Liphium/station/spacestation/util"
)

// Action: setup
func setup(c *pipeshandler.Context, action struct {
	Data  string  `json:"data"`
	Color float64 `json:"color"`
}) pipes.Event {

	// Generate new connection
	connection := caching.EmptyConnection(c.Client.ID, c.Client.Session)

	// Insert data
	if !caching.SetMemberData(c.Client.Session, c.Client.ID, connection.ClientID, action.Data) {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	// Send the update to all members in the room
	if !SendRoomData(c.Client.Session) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	// Have the guy join the table
	msg := caching.JoinTable(c.Client.Session, c.Client.ID, action.Color)
	if msg != nil {
		util.Log.Println("Couldn't join table of room", c.Client.Session, ":", msg[localization.DefaultLocale])
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"id":      connection.ClientID,
		"key":     connection.KeyBase64(),
	})
}
