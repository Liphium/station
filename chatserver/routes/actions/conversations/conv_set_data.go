package conversation_actions

import (
	"github.com/Liphium/station/chatserver/database"
	message_actions "github.com/Liphium/station/chatserver/routes/actions/messages"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type changeDataRequest struct {
	Version int64  `json:"version"`
	Data    string `json:"data"`
}

// Action: conv_set_data
func HandleSetData(c *fiber.Ctx, token database.ConversationToken, action changeDataRequest) error {

	// Make sure the person has at least the moderator rank
	if token.Rank < database.RankModerator {
		return integration.FailedRequest(c, localization.ErrorMemberNoPermission, nil)
	}

	// Edit the conversation
	// The version here makes sure that no other person is editing the conversation at the same time.
	// The query will fail on updates at the same time, but since this is only a protection for the worst
	// case scenario, we should be fine without a specific error here. It's gonna be fine.. hopefully :)
	if err := database.DBConn.Where("id = ? AND type != ? AND version = ?", token.Conversation, database.ConvTypePrivateMessage, action.Version).Updates(&database.Conversation{
		Version: action.Version + 1,
		Data:    action.Data,
	}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send a system to everyone to tell them about the change of the data
	if err := message_actions.SendSystemMessage(token.Conversation, "", message_actions.ConversationEdited, []string{
		message_actions.AttachAccount(token.Data),
	}); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
