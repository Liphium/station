package conversation

import (
	"slices"
	"strings"
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/handler/account"
	conversation_actions "github.com/Liphium/station/chatserver/routes/actions/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/gofiber/fiber/v2"
)

type conversationSubscribeResponse struct {
	Success bool                                             `json:"success"`
	Info    map[string]conversation_actions.ConversationInfo `json:"info"`
	Missing []string                                         `json:"missing"`
	Node    string                                           `json:"node"`
}

// Action: conv_sub
func subscribe(c *pipeshandler.Context, action struct {
	Tokens   []conversations.SentConversationToken `json:"tokens"`
	Status   string                                `json:"status"`
	SyncDate int64                                 `json:"sync"` // Time of last sent message for message sync
	Data     string                                `json:"data"`
}) pipes.Event {

	// Filter out all the remote tokens and register adapters
	localTokens := []conversations.SentConversationToken{}
	remoteTokens := map[string][]conversations.SentConversationToken{}
	for _, token := range action.Tokens {

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
				caching.CSNode.RemoveAdapterWS("s-" + token.ID)
				caching.CSInstance.Disconnect(c.Client.ID, c.Client.Session)
			},
		})

		// Extract the address
		args := strings.Split(token.ID, "@")
		if len(args) != 2 {
			continue
		}

		// Check if a remote subscription should be registered
		if args[1] != integration.Domain {

			// Add the token to the remote tokens for that instance
			if remoteTokens[args[1]] == nil {
				remoteTokens[args[1]] = []conversations.SentConversationToken{token}
			} else {
				remoteTokens[args[1]] = append(remoteTokens[args[1]], token)
			}
		} else {
			localTokens = append(localTokens, token)
		}
	}

	// Validate the tokens
	conversationTokens, missingTokens, conversationIds, err := caching.ValidateTokens(&localTokens)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Get the conversation info
	convInfo, err := conversation_actions.GetConversationInfo(conversationTokens)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	// Subscribe to all remote tokens
	var serversWithError []string = []string{}
	enabled, err := integration.GetBoolSetting(caching.CSNode, integration.SettingDecentralizationEnabled)
	if err != nil || !enabled {
		for server := range remoteTokens {
			serversWithError = append(serversWithError, server)
		}
	}

	if enabled {
		for server, tokens := range remoteTokens {
			res, err := action_helpers.SendRemoteActionGeneric[conversationSubscribeResponse](server, "conv_subscribe", fiber.Map{
				"tokens": tokens,
				"status": action.Status,
				"data":   action.Data,
				"node":   util.OwnPath,
			})

			// Check if there was an error, if so, tell the client
			if err != nil || !res.Success {
				serversWithError = append(serversWithError, server)
				continue
			}

			// Add the conversation info from the remote server
			// This could technically be vulnerable to an attack where a remote node could
			// artificially increment the notification count, mess with the read dates or
			// make the client re-fetch the conversation version (just why?). To me, this isn't
			// of importance and because this would need a lot of code changes to fix, I'll
			// just leave this reminder here. If anyone finds this in the future, have
			// fun exploiting this! :D
			// And if you find anything actually serious, well, you know who to blame this on.
			for conv, info := range res.Answer.Info {
				convInfo[conv] = info
			}

			// Add the missing tokens
			// Make sure remote nodes can't delete tokens they don't have access to (important security fix)
			res.Answer.Missing = slices.DeleteFunc(res.Answer.Missing, func(element string) bool {
				return !slices.ContainsFunc(tokens, func(token conversations.SentConversationToken) bool {
					return token.ID == element
				})
			})
			missingTokens = append(missingTokens, res.Answer.Missing...)
		}
	}

	// Send the status to everyone in a goroutine
	go func() {

		// Grab all the members
		members, err := caching.LoadMembersArray(conversationIds)
		if err != nil {
			return
		}

		// Send the status to all the conversations
		for _, token := range conversationTokens {

			// Make sure it's only send to private conversations
			if len(members[token.Conversation]) > 2 {
				continue
			}

			// Get the other member to send the status to
			var otherMember caching.StoredMember
			for _, member := range members[token.Conversation] {
				if member.TokenID != token.ID {
					otherMember = member
				}
			}

			// Send the status event
			caching.SendEventToMembers([]caching.StoredMember{otherMember}, account.StatusEvent(action.Status, action.Data, token.Conversation, token.ID, ""))
		}
	}()

	// Synchronize messages for the local tokens
	go func() {
		// Wait for the client to receive the response
		time.Sleep(1 * time.Second)

		// Go through local tokens to add them to the message sync queue (if desired)
		if action.SyncDate != -1 {
			for _, token := range conversationTokens {
				if token.Activated {
					if err := caching.AddSyncToQueue(caching.SyncData{
						TokenID:      token.ID,
						Conversation: token.Conversation,
						Since:        action.SyncDate,
					}); err != nil {
						util.Log.Println("error completing message sync for ", token.ID, ":", err)
					}
				}
			}
		}
	}()

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"info":    convInfo,
		"error":   serversWithError,
		"missing": missingTokens,
	})
}
