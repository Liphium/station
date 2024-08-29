package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

type statusRespondAction struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

// Action: st_res
func respondToStatus(c *pipeshandler.Context, action statusRespondAction) pipes.Event {

	// Get from cache
	convToken, err := caching.ValidateToken(action.ID, action.Token)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.InvalidRequest, err)
	}

	// Check if this is a valid conversation
	members, err := caching.LoadMembers(convToken.Conversation)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Make sure it's a private conversation
	if len(members) > 2 {
		return pipeshandler.ErrorResponse(c, localization.InvalidRequest, nil)
	}

	// Get the other member to send the status to
	var otherMember caching.StoredMember
	for _, member := range members {
		if member.TokenID != convToken.ID {
			otherMember = member
		}
	}

	// Send the event
	if err := caching.SendEventToMembers([]caching.StoredMember{otherMember}, StatusEvent(action.Status, action.Data, convToken.Conversation, convToken.ID, ":a")); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.SuccessResponse(c)
}
