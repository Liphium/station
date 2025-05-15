package conversation_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/handler/account"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_st_res
func HandleStatusResponse(c *fiber.Ctx, token database.ConversationToken, action struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}) error {

	// Make sure the token is activated
	if !token.Activated {
		return integration.InvalidRequest(c, "not activated token")
	}

	// Check if this is a valid conversation
	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Make sure it's a private conversation
	if len(members) > 2 {
		return integration.InvalidRequest(c, "conversation isn't a private conversation")
	}

	// Get the other member to send the status to
	var otherMember caching.StoredMember
	for _, member := range members {
		if member.TokenID != token.ID {
			otherMember = member
		}
	}

	// Send the event
	if err := caching.SendEventToMembers([]caching.StoredMember{otherMember}, account.StatusEvent(action.Status, action.Data, token.Conversation, token.ID, ":a")); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
