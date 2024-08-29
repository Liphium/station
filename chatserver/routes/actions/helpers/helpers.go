package action_helpers

import (
	"strings"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// A generic type for a request to any conversation remote action
type ConversationActionRequest[T any] struct {
	Token conversations.SentConversationToken `json:"token"`
	Data  T                                   `json:"data"`
}

// A generic type for any action handler function
type ConversationActionHandlerFunc[T any] func(*fiber.Ctx, conversations.ConversationToken, T) error

// A generic type for any action handler function
type ActionHandlerFunc[T any] func(*fiber.Ctx, T) error

// Generic struct to wrap the json with any additional data for an action
type RemoteActionRequest[T any] struct {
	ID     string `json:"id"`
	Token  string `json:"token"`
	Sender string `json:"sender"`
	Data   T      `json:"data"`
}

// This deletes all data related to a conversation
func DeleteConversation(id string) error {
	if err := database.DBConn.Where("conversation = ?", id).Delete(&conversations.Message{}).Error; err != nil {
		return err
	}
	if err := database.DBConn.Where("conversation = ?", id).Delete(&conversations.ConversationToken{}).Error; err != nil {
		return err
	}
	if err := database.DBConn.Where("id = ?", id).Delete(&conversations.Conversation{}).Error; err != nil {
		return err
	}
	return nil
}

// This increments the version of the conversation by one in a transaction.
// Will also save the conversation.
func IncrementConversationVersion(conversation conversations.Conversation) error {

	// Increment the version in a transaction
	err := database.DBConn.Transaction(func(tx *gorm.DB) error {

		// Get the current version (in case it has changed)
		var currentVersion int64
		if err := tx.Model(&conversations.Conversation{}).Select("version").Where("id = ?", conversation.ID).Take(&currentVersion).Error; err != nil {
			database.DBConn.Rollback()
			return err
		}

		// Update the conversation
		conversation.Version = currentVersion + 1

		// Save the conversation
		if err := tx.Save(&conversation).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// Create a normal endpoint from an conversation action handler
func CreateConversationEndpoint[T any](handler ConversationActionHandlerFunc[T], action string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// Parse the request
		var req ConversationActionRequest[T]
		if err := integration.BodyParser(c, &req); err != nil {
			return integration.InvalidRequest(c, "request was invalid")
		}

		// Parse the conversation to extract the address
		args := strings.Split(req.Token.ID, "@")
		if len(args) != 2 {
			return integration.InvalidRequest(c, "conversation id is invalid")
		}

		// If the address isn't the current instance, send a remote action
		if args[1] != integration.BasePath {

			// Send a remote action to the other instance
			res, err := integration.PostRequestBackendServer(args[1], "/node/actions/"+action, fiber.Map{
				"app_tag": integration.AppTagChatNode,
				"sender":  caching.CSNode.SL,
				"action":  action,
				"data":    req,
			})
			if err != nil {
				return integration.FailedRequest(c, localization.ErrorServer, err)
			}

			// Check if the request was successful
			if !res["success"].(bool) {
				return integration.FailedRequest(c, localization.ErrorNode, err)
			}

			// Return the response to the client
			return integration.ReturnJSON(c, res["answer"])
		}

		// Validate the token
		token, err := caching.ValidateToken(req.Token.ID, req.Token.Token)
		if err != nil {
			return integration.InvalidRequest(c, "conversation token was valid")
		}

		// Let the action handle the request
		return handler(c, token, req.Data)
	}
}
