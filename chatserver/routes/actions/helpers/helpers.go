package action_helpers

import (
	"errors"
	"strings"
	"sync"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
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

		// Parse the conversation token id to extract the address
		args := strings.Split(req.Token.ID, "@")
		if len(args) != 2 {
			return integration.InvalidRequest(c, "conversation id is invalid")
		}

		// Check if the connection is safe (or if unsafe is allowed)
		if strings.HasPrefix(strings.TrimSpace(args[1]), "http://") && !util.AllowUnsafe {
			return integration.FailedRequest(c, localization.ErrorNoUnsafeConnections, errors.New("unsafe requests aren't allowed"))
		}

		// If the address isn't the current instance, send a remote action
		if args[1] != integration.Domain {

			// Send a conversation aciton to the other instance
			res, err := SendConversationAction(action, req.Token, req.Data)
			if err != nil {

				// Check if it's an error that happend on the other server
				if strings.Contains(err.Error(), "other server error:") {
					return integration.FailedRequest(c, localization.ErrorOtherServer, err)
				}

				return integration.FailedRequest(c, localization.ErrorServer, err)
			}

			// Return the response to the client
			return integration.ReturnJSON(c, res)
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

// Send a conversation action using a conversation token
func SendConversationAction(action string, token conversations.SentConversationToken, data interface{}) (map[string]interface{}, error) {

	// Get the address of the chat node
	obj, valid := TokenMap.Load(token.ID)
	if !valid {

		// Extract the address of the main backend from the conversation token
		args := strings.Split(token.ID, "@")
		if len(args) != 2 {
			return nil, errors.New("address of conversation token couldn't be parsed")
		}

		// Negotiate with the chat server to get its address
		if err := negotiate(args[1], token.ID, token.Token); err != nil {
			return nil, errors.New("other server error: " + err.Error())
		}
		obj, valid = TokenMap.Load(token.ID)
	}

	// Make sure the token exists after negotiation
	if !valid {
		return nil, errors.New("token couldn't be found")
	}
	node := obj.(*TokenData)

	// Send the action
	var res map[string]interface{}
	var err error
	if res, err = integration.PostRequestTC(node.Node, "/conv_actions/"+action, fiber.Map{
		"token": token,
		"data":  data,
	}); err != nil {
		return nil, errors.New("other server error: " + err.Error())
	}

	return res, nil
}

type TokenData struct {
	Token string // The actual token
	Node  string // The address of the node subscribed to
}

// Token -> *TokenData
var TokenMap *sync.Map = &sync.Map{}

// Send a negotiation offer to any node
func negotiate(server string, id string, token string) error {

	// Make sure this thing can't be crashed
	defer func() {
		if err := recover(); err != nil {
			util.Log.Println("something went seriously wrong with negotiation: ", err)
		}
	}()

	// Send the remote action
	res, err := SendRemoteAction(server, "negotiate", fiber.Map{
		"id":    id,
		"token": token,
		"node":  util.OwnPath,
	})

	// Check if there was some kind of error
	if err != nil {
		return err
	}
	if !res["success"].(bool) {
		return errors.New("remote action couldn't be sent: " + res["error"].(string))
	}

	// Extract the answer
	answer := res["answer"].(map[string]interface{})
	if !answer["success"].(bool) {
		return errors.New("negotiation was declined: " + answer["error"].(string))
	}

	// Store the data from the request in the token map
	StoreToken(conversations.SentConversationToken{ID: id, Token: token}, answer["node"].(string))

	return nil
}

// Store any conversation token in the token map (make sure it can be reached from outside)
func StoreToken(token conversations.SentConversationToken, node string) {
	TokenMap.Store(token.ID, &TokenData{
		Token: token.Token,
		Node:  node,
	})
}

// Sends a remote action to any server
func SendRemoteAction(server string, action string, data interface{}) (map[string]interface{}, error) {
	return integration.PostRequestBackendServer(server, "/node/actions/send", fiber.Map{
		"app_tag": integration.AppTagChatNode,
		"sender":  integration.BasePath,
		"action":  action,
		"data":    data,
	})
}

type remoteActionResponse[T any] struct {
	Success bool `json:"success"`
	Answer  T    `json:"answer"`
}

// Sends a remote action to any server using a generic response type
func SendRemoteActionGeneric[T any](server string, action string, data interface{}) (remoteActionResponse[T], error) {
	return integration.PostRequestBackendServerGeneric[remoteActionResponse[T]](server, "/node/actions/send", fiber.Map{
		"app_tag": integration.AppTagChatNode,
		"sender":  integration.BasePath,
		"action":  action,
		"data":    data,
	})
}
