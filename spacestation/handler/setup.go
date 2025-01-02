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
	Data      string `json:"data"`
	Signature string `json:"signature"`
}) pipes.Event {

	// Insert data
	if !caching.SetMemberData(c.Client.Session, c.Client.ID, action.Data, action.Signature) {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	// Send the update to all members in the room
	if !SendRoomData(c.Client.Session) {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	// Have the guy join the table
	msg := caching.JoinTable(c.Client.Session, c.Client.ID)
	if msg != nil {
		util.Log.Println("Couldn't join table of room", c.Client.Session, ":", msg[localization.DefaultLocale])
		return pipeshandler.ErrorResponse(c, msg, nil)
	}

	// Send the guy all the warps
	caching.InitializeWarps(c.Client)

	return pipeshandler.SuccessResponse(c)
}
