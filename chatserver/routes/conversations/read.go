package conversation_routes

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Public so it can be unit tested (in the future ig)
type ReadConversationRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (r *ReadConversationRequest) Validate() bool {
	return len(r.ID) > 0 && len(r.Token) > 0 && len(r.Token) == util.ConversationTokenLength
}

// Route: /conversations/read
func read(c *fiber.Ctx) error {

	var req ReadConversationRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, err.Error())
	}

	// Validate request
	if !req.Validate() {
		util.Log.Println(len(req.Token))
		return integration.InvalidRequest(c, "request is invalid")
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "token is invalid")
	}

	// Update read state
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ?", token.ID).Update("last_read", time.Now().UnixMilli()).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	token.LastRead = time.Now().UnixMilli()

	return integration.SuccessfulRequest(c)
}
