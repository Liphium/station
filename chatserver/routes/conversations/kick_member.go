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

type kickMemberRequest struct {
	Id     string `json:"id"`
	Token  string `json:"token"`
	Target string `json:"target"`
}

// Route: /conversations/kick_member
func kickMember(c *fiber.Ctx) error {

	var req kickMemberRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	if req.Id == req.Target {
		return integration.InvalidRequest(c, "same token")
	}

	token, err := caching.ValidateToken(req.Id, req.Token)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	targetToken, err := caching.GetToken(req.Target)
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
	if err := incrementConversationVersion(conversation); err != nil {
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
