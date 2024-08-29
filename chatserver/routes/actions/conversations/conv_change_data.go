package conversation_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_actions "github.com/Liphium/station/chatserver/routes/actions/messages"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

type ChangeDataAction struct {
	Id    string `json:"id"`
	Token string `json:"token"`
	Data  string `json:"data"`
}

// Action: conv_change_data
func HandleChangeData(c *fiber.Ctx, token conversations.ConversationToken, data string) error {

	// Check if the form is valid
	if len(data) > util.MaxConversationDataLength {
		return integration.FailedRequest(c, localization.GroupDataTooLong, nil)
	}

	// Get the conversation
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Take(&conversation).Error; err != nil {
		database.DBConn.Rollback()
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if it is a group
	if conversation.Type != conversations.TypeGroup {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Update the data
	conversation.Data = data

	// Increment the version by one to let everyone know
	if err := action_helpers.IncrementConversationVersion(conversation); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send a system to everyone to tell them about the change of the data
	if err := message_actions.SendSystemMessage(token.Conversation, message_actions.ConversationEdited, []string{
		message_actions.AttachAccount(token.Data),
	}); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return success
	return integration.SuccessfulRequest(c)
}
