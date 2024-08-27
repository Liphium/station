package conversation_actions

import (
	"fmt"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type KickMemberAction struct {
	Id     string `json:"id"`
	Token  string `json:"token"`
	Target string `json:"target"`
}

// Action: conv_kick
func HandleKick(c *fiber.Ctx, action KickMemberAction) error {

	// Make sure you can't kick yourself
	if action.Id == action.Target {
		return integration.InvalidRequest(c, "same token")
	}

	// Validate the token and get all the tokens
	token, err := caching.ValidateToken(action.Id, action.Token)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	targetToken, err := caching.GetToken(action.Target)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if conversation is group
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error; err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in database: %s", err.Error()))
	}

	if conversation.Type != conversations.TypeGroup {
		return integration.FailedRequest(c, "no.group", nil)
	}

	// Check if the token has the permission
	if token.Rank <= targetToken.Rank {
		return integration.FailedRequest(c, localization.KickNoPermission, nil)
	}

	// Increment the version by one to save the modification
	if err := action_helpers.IncrementConversationVersion(conversation); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Delete from the database
	if err := database.DBConn.Delete(&targetToken).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupMemberKick, []string{message_routes.AttachAccount(token.Data), message_routes.AttachAccount(targetToken.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	err = message_routes.SendNotStoredSystemMessage(token.Conversation, message_routes.ConversationKick, []string{message_routes.AttachAccount(targetToken.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Unsubscribe from stuff
	caching.DeleteToken(targetToken.ID, targetToken.Token)

	return integration.SuccessfulRequest(c)
}
