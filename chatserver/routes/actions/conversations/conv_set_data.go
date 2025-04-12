package conversation_actions

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type changeDataRequest struct {
	Version int64  `json:"version"`
	Data    string `json:"data"`
}

// Action: conv_set_data
func HandleSetData(c *fiber.Ctx, token conversations.ConversationToken, action changeDataRequest) error {

	// Edit the conversation
	// The version here makes sure that no other person is editing the conversation at the same time.
	// The query will fail on updates at the same time, but since this is only a protection for the worst
	// case scenario, we should be fine without a specific error here. It's gonna be fine.. hopefully :)
	if err := database.DBConn.Where("id = ? AND version = ?", token.Conversation, action.Version).Updates(&conversations.Conversation{
		Version: action.Version + 1,
		Data:    action.Data,
	}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
