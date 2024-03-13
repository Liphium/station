package conversation_routes

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	message_routes "github.com/Liphium/station/chatserver/routes/conversations/message"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Public so it can be unit tested (in the future ig)
type ActivateConversationRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (r *ActivateConversationRequest) Validate() bool {
	return len(r.ID) > 0 && len(r.Token) > 0 && len(r.Token) == util.ConversationTokenLength
}

type returnableMember struct {
	ID   string `json:"id"`
	Rank uint   `json:"rank"`
	Data string `json:"data"`
}

// Route: /conversations/activate
func activate(c *fiber.Ctx) error {

	var req ActivateConversationRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, err.Error())
	}

	// Validate request
	if !req.Validate() {
		util.Log.Println(len(req.Token))
		return integration.InvalidRequest(c, "request is invalid")
	}

	// Activate conversation
	var token conversations.ConversationToken
	if database.DBConn.Where("id = ? AND token = ?", req.ID, req.Token).First(&token).Error != nil {
		return integration.FailedRequest(c, "invalid.token", nil)
	}

	if token.Activated {
		return integration.FailedRequest(c, "already.active", nil)
	}

	// Activate token
	token.Activated = true
	token.Token = util.GenerateToken(util.ConversationTokenLength)

	if err := database.DBConn.Save(&token).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Update token
	err := caching.UpdateToken(token)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return all data
	var tokens []conversations.ConversationToken
	if err := database.DBConn.Where(&conversations.ConversationToken{Conversation: token.Conversation}).Find(&tokens).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	var members []returnableMember
	for _, token := range tokens {
		members = append(members, returnableMember{
			ID:   token.ID,
			Rank: token.Rank,
			Data: token.Data,
		})
	}

	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Take(&conversation).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if conversation.Type == conversations.TypeGroup {
		err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupMemberJoin, []string{message_routes.AttachAccount(token.Data)})
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"type":    conversation.Type,
		"data":    conversation.Data,
		"token":   token.Token,
		"members": members,
	})
}
