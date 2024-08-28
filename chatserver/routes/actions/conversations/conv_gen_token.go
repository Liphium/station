package conversation_actions

import (
	"fmt"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_gen_token
func HandleGenerateToken(c *fiber.Ctx, token conversations.ConversationToken, data string) error {

	// Check if conversation is group
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error; err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in database: %s", err.Error()))
	}

	if conversation.Type != conversations.TypeGroup {
		return integration.FailedRequest(c, localization.GroupInvalidType, nil)
	}

	// Check requirements for a new token
	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if len(members) >= 100 {
		return integration.FailedRequest(c, localization.GroupMemberLimit, nil)
	}

	// Increment the version by one to save the modification
	if err := action_helpers.IncrementConversationVersion(conversation); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Generate a new token
	generated := conversations.ConversationToken{
		ID:           util.GenerateToken(util.ConversationTokenIDLength),
		Token:        util.GenerateToken(util.ConversationTokenLength),
		Activated:    false,
		Conversation: token.Conversation,
		Rank:         conversations.RankUser,
		Data:         data,
	}

	if err := database.DBConn.Create(&generated).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send a system message to let everyone know
	err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupMemberInvite, []string{message_routes.AttachAccount(token.Data), message_routes.AttachAccount(generated.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"id":      generated.ID,
		"token":   generated.Token,
	})
}
