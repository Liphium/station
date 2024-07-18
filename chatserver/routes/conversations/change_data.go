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

type changeDataRequest struct {
	Id    string `json:"id"`
	Token string `json:"token"`
	Data  string `json:"data"`
}

// Route: /conversations/change_data
func changeData(c *fiber.Ctx) error {

	// Parse the request
	var req changeDataRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Check if the form is valid
	if len(req.Data) > util.MaxConversationDataLength {
		return integration.FailedRequest(c, localization.GroupDataTooLong, nil)
	}

	// Validate the token
	token, err := caching.ValidateToken(req.Id, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "invalid token")
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
	conversation.Data = req.Data

	// Increment the version by one to let everyone know
	if err := incrementConversationVersion(conversation); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send a system to everyone to tell them about the change of the data
	if err := message_routes.SendSystemMessage(token.Conversation, message_routes.ConversationEdited, []string{
		message_routes.AttachAccount(token.Data),
	}); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Return success
	return integration.SuccessfulRequest(c)
}
