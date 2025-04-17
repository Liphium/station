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

type KickMemberAction struct {
	Id     string `json:"id"`
	Token  string `json:"token"`
	Target string `json:"target"`
}

// Action: conv_kick
func HandleKick(c *fiber.Ctx, token database.ConversationToken, target string) error {

	// Make sure you can't kick yourself
	if token.ID == target {
		return integration.InvalidRequest(c, "same token")
	}

	// Get the token of the target
	targetToken, err := caching.GetToken(target)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Make sure the conversation isn't a private message
	var conversation database.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error; err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in database: %s", err.Error()))
	}

	if conversation.Type == database.ConvTypePrivateMessage {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Check if the token has the permission
	if token.Rank <= targetToken.Rank {
		return integration.FailedRequest(c, localization.ErrorKickNoPermission, nil)
	}

	// Increment the version by one to save the modification
	if err := action_helpers.IncrementConversationVersion(conversation); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete from the database
	if err := database.DBConn.Delete(&targetToken).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	err = message_actions.SendSystemMessage(token.Conversation, "", message_actions.GroupMemberKick, []string{message_actions.AttachAccount(token.Data), message_actions.AttachAccount(targetToken.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	err = message_actions.SendNotStoredSystemMessage(token.Conversation, "", message_actions.ConversationKick, []string{message_actions.AttachAccount(targetToken.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Unsubscribe from stuff
	caching.DeleteToken(targetToken.ID, targetToken.Token)

	return integration.SuccessfulRequest(c)
}
