package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/database/fetching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

type sendStatusAction struct {
	Tokens []conversations.SentConversationToken `json:"tokens"`
	Status string                                `json:"status"`
	Data   string                                `json:"data"`
}

// Action: st_send
func sendStatus(c *pipeshandler.Context, action sendStatusAction) pipes.Event {

	// Save in database
	if err := database.DBConn.Model(&fetching.Status{}).Where("id = ?", c.Client.ID).Update("data", action.Data).Error; err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Validate all the tokens
	conversationTokens, _, tokenIds, err := caching.ValidateTokens(&action.Tokens)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Grab all the members
	members, err := caching.LoadMembersArray(tokenIds)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	for _, token := range conversationTokens {

		// Make sure the conversation is a private one
		if len(members[token.Conversation]) > 2 {
			continue
		}

		// Get the other member to send the status to
		var otherMember caching.StoredMember
		for _, member := range members[token.Conversation] {
			if member.TokenID != token.ID {
				otherMember = member
			}
		}

		// Send the status event
		caching.SendEventToMembers([]caching.StoredMember{otherMember}, StatusEvent(action.Status, action.Data, token.Conversation, token.ID, ""))
	}

	// Send the status to other devices
	caching.CSNode.SendClient(c.Client.ID, StatusEvent(action.Status, action.Data, "", c.Client.ID, ":o"))

	return pipeshandler.SuccessResponse(c)
}

func StatusEvent(st string, data string, conversation string, ownToken string, suffix string) pipes.Event {
	return pipes.Event{
		Name: "acc_st" + suffix,
		Data: map[string]interface{}{
			"c":  conversation,
			"o":  ownToken,
			"st": st,
			"d":  data,
		},
	}
}
