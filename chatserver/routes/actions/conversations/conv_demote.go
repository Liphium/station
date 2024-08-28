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

// Action: conv_demote
func HandleDemoteToken(c *fiber.Ctx, token conversations.ConversationToken, user string) error {

	// Check if conversation is group
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error; err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in database: %s", err.Error()))
	}

	if conversation.Type != conversations.TypeGroup {
		return integration.FailedRequest(c, "no.group", nil)
	}

	if token.Rank == conversations.RankUser {
		return integration.InvalidRequest(c, "user doesn't have the required rank")
	}

	// Get the token of the other user
	userToken, err := caching.GetToken(user)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("specified user doesn't exist: %s", err.Error()))
	}

	if userToken.Conversation != token.Conversation {
		return integration.InvalidRequest(c, "conversations don't match")
	}

	// Get rank to promote (check permissions)
	rankToDemote := userToken.Rank - 1
	if userToken.Rank > token.Rank {
		return integration.InvalidRequest(c, "no permission")
	}

	// Increment the version by one to save the modification
	if err := action_helpers.IncrementConversationVersion(conversation); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Demote the person in the database
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ? AND conversation = ?", userToken.ID, userToken.Conversation).Update("rank", rankToDemote).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Set the rank for the system message
	prevRank := userToken.Rank
	userToken.Rank = rankToDemote

	// Send a system message to let everyone know about the rank change
	err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupRankChange, []string{fmt.Sprintf("%d", prevRank), fmt.Sprintf("%d", userToken.Rank),
		message_routes.AttachAccount(userToken.Data), message_routes.AttachAccount(token.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
