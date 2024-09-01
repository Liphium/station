package conversation

import (
	"slices"
	"strings"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database/conversations"
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/gofiber/fiber/v2"
)

type conversationSubscribeResponse struct {
	Success bool                                             `json:"success"`
	Info    map[string]conversation_actions.ConversationInfo `json:"info"`
	Missing []string                                         `json:"missing"`
}

// Action: conv_sub
func subscribe(c *pipeshandler.Context, action struct {
	Tokens []conversations.SentConversationToken `json:"tokens"`
	Status string                                `json:"status"`
}) pipes.Event {

	// Validate the tokens
	conversationTokens, missingTokens, _, err := caching.ValidateTokens(&action.Tokens)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Get the conversation info
	convInfo, err := conversation_actions.GetConversationInfo(conversationTokens)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Register all the adapters
	remoteTokens := map[string][]conversations.SentConversationToken{}
	for _, token := range conversationTokens {
		if token.Activated {

			// Extract the address
			args := strings.Split(token.ID, "@")
			if len(args) != 2 {
				continue
			}

			// Register adapter for the subscription
			caching.CSNode.AdaptWS(pipes.Adapter{
				ID: "s-" + token.ID,
				Receive: func(context *pipes.Context) error {
					client := *c.Client
					util.Log.Println(context.Adapter.ID, token.Token, client.ID)
					err := caching.CSNode.SendClient(c.Client.ID, *context.Event)
					if err != nil {
						util.Log.Println("COULDN'T SEND:", err.Error())
					}
					return err
				},

				// Remove the adapter if there was an error (and disconnect the user)
				OnError: func(err error) {
					caching.CSNode.RemoveAdapterWS("s-" + token.Token)
					caching.CSInstance.Disconnect(c.Client.ID, c.Client.Session)
				},
			})

			// Check if a remote subscription should be registered
			if args[1] != integration.Domain {
				sentToken := token.ToSent()

				// Add the token to the remote tokens for that instance
				if remoteTokens[args[1]] == nil {
					remoteTokens[args[1]] = []conversations.SentConversationToken{sentToken}
				} else {
					remoteTokens[args[1]] = append(remoteTokens[args[1]], sentToken)
				}
			}
		}
	}

	// Subscribe to all remote tokens
	var serversWithError []string = []string{}
	for server, tokens := range remoteTokens {
		res, err := action_helpers.SendRemoteActionGeneric[conversationSubscribeResponse](server, "conv_subscribe", fiber.Map{
			"tokens": tokens,
			"status": action.Status,
		})

		// Check if there was an error, if so, tell the client
		if err != nil {
			serversWithError = append(serversWithError, server)
		}

		// Add the conversation info from the remote server
		// This loop is coded this way for security reasons (so the other server couldn't delete a conversation that is not on it)
		for _, token := range tokens {
			convInfo[token.ID] = res.Info[token.ID]
		}

		// Add the missing tokens
		res.Missing = slices.DeleteFunc(res.Missing, func(element string) bool {
			return !slices.ContainsFunc(tokens, func(token conversations.SentConversationToken) bool {
				return token.ID == element
			})
		})
		missingTokens = append(missingTokens, res.Missing...)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"info":    convInfo,
		"error":   serversWithError,
		"missing": missingTokens,
	})
}
