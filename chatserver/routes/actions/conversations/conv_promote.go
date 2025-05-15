package conversation_actions

import (
	"fmt"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_actions "github.com/Liphium/station/chatserver/routes/actions/messages"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_promote
func HandlePromoteToken(c *fiber.Ctx, token database.ConversationToken, user string) error {

	// Make sure the token is activated
	if !token.Activated {
		return integration.InvalidRequest(c, "not activated token")
	}

	// Make sure the conversation is not a private message
	var conversation database.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error; err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in db: %s", err.Error()))
	}
	if conversation.Type == database.ConvTypePrivateMessage {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	if token.Rank == database.RankUser {
		return integration.InvalidRequest(c, "no permission")
	}

	// Get the token of the other user
	userToken, err := caching.GetToken(user)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't get user token: %s", err.Error()))
	}

	if userToken.Conversation != token.Conversation {
		return integration.InvalidRequest(c, "conversations don't match")
	}

	// Get rank to promote (check permissions)
	rankToPromote := userToken.Rank + 1
	if rankToPromote > token.Rank {
		return integration.InvalidRequest(c, "no permission for promotion")
	}

	// Increment the version by one to save the modification
	if err := action_helpers.IncrementConversationVersion(conversation); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Update the rank in the database
	if err := database.DBConn.Model(&database.ConversationToken{}).Where("id = ? AND conversation = ?", userToken.ID, userToken.Conversation).Update("rank", rankToPromote).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Change it for the system message
	prevRank := userToken.Rank
	userToken.Rank = rankToPromote

	// Send a system message to let all members know about the rank change
	err = message_actions.SendSystemMessage(token.Conversation, "", message_actions.GroupRankChange, []string{fmt.Sprintf("%d", prevRank), fmt.Sprintf("%d", userToken.Rank),
		message_actions.AttachAccount(userToken.Data), message_actions.AttachAccount(token.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
