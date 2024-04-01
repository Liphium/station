package conversation_routes

import (
	"fmt"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Route: /conversations/demote_token
func demoteToken(c *fiber.Ctx) error {

	var req promoteTokenRequest
	if integration.BodyParser(c, &req) != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("invalid token: %s", err.Error()))
	}

	// Check if conversation is group
	var conversation conversations.Conversation
	if database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in database: %s", err.Error()))
	}

	if conversation.Type != conversations.TypeGroup {
		return integration.FailedRequest(c, "no.group", nil)
	}

	if token.Rank == conversations.RankUser {
		return integration.InvalidRequest(c, "user doesn't have the required rank")
	}

	userToken, err := caching.GetToken(req.User)
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

	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ? AND conversation = ?", userToken.ID, userToken.Conversation).Update("rank", rankToDemote).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	prevRank := userToken.Rank
	userToken.Rank = rankToDemote

	err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupRankChange, []string{fmt.Sprintf("%d", prevRank), fmt.Sprintf("%d", userToken.Rank),
		message_routes.AttachAccount(userToken.Data), message_routes.AttachAccount(token.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
