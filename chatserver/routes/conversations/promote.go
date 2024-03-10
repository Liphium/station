package conversation_routes

import (
	"fmt"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/integration"
	"github.com/gofiber/fiber/v2"
)

type promoteTokenRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	User  string `json:"user"`
}

// Route: /conversations/promote_token
func promoteToken(c *fiber.Ctx) error {

	var req promoteTokenRequest
	if integration.BodyParser(c, &req) != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("invalid conversation token: %s", err.Error()))
	}

	// Check if conversation is group
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Find(&conversation).Error; err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't find conversation in db: %s", err.Error()))
	}

	if conversation.Type != conversations.TypeGroup {
		return integration.FailedRequest(c, "no.group", nil)
	}

	if token.Rank == conversations.RankUser {
		return integration.InvalidRequest(c, "no permission")
	}

	userToken, err := caching.GetToken(req.User)
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

	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ? AND conversation = ?", userToken.ID, userToken.Conversation).Update("rank", rankToPromote).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	prevRank := userToken.Rank
	userToken.Rank = rankToPromote
	err = caching.UpdateToken(userToken)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupRankChange, []string{fmt.Sprintf("%d", prevRank), fmt.Sprintf("%d", userToken.Rank),
		message_routes.AttachAccount(userToken.Data), message_routes.AttachAccount(token.Data)})
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
