package conversation

import (
	"slices"
	"strings"
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
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
	Tokens []database.SentConversationToken `json:"tokens"`
	Status string                           `json:"status"`
	Data   string                           `json:"data"`
}) pipes.Event {

	// Filter out all the remote tokens and register adapters
	localTokens := []database.SentConversationToken{}
	remoteTokens := map[string][]database.SentConversationToken{}
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
				remoteTokens[args[1]] = []database.SentConversationToken{token}
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
	enabled, err := integration.GetBoolSetting(caching.CSNode, integration.SettingDecentralizationEnabled)
	if enabled && err == nil {
		for server, tokens := range remoteTokens {

			// These subscriptions will be delivered to the client later, due to potential timeouts taking way too long
			// and making the client feel really slow.
			go func() {
				res, err := action_helpers.SendRemoteActionGeneric[conversationSubscribeResponse](server, "conv_subscribe", fiber.Map{
					"tokens": tokens,
					"status": action.Status,
					"data":   action.Data,
					"node":   util.OwnPath,
				})

				// Check if there was an error, if so, tell the client
				if err != nil || !res.Success {

					// Send an error for this server
					caching.CSInstance.SendEventToOne(c.Client, pipes.Event{
						Name: "conv_sub:late",
						Data: map[string]interface{}{
							"server": server,
							"error":  true,
						},
					})
					return
				}

				// Make sure remote nodes can't delete tokens they don't have access to (important security fix)
				res.Answer.Missing = slices.DeleteFunc(res.Answer.Missing, func(element string) bool {
					return !slices.ContainsFunc(tokens, func(token database.SentConversationToken) bool {
						return token.ID == element
					})
				})

				// Send all the information for this server
				caching.CSInstance.SendEventToOne(c.Client, pipes.Event{
					Name: "conv_sub:late",
					Data: map[string]interface{}{
						"server":  server,
						"error":   false,
						"missing": res.Answer.Missing,
						"info":    res.Answer.Info,
					},
				})
			}()
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

		// Go through local tokens to add them to the message sync queue
		for _, token := range conversationTokens {
			if token.Activated && token.LastSync != -1 {
				if err := caching.AddSyncToQueue(caching.SyncData{
					TokenID:      token.ID,
					Conversation: token.Conversation,
					Since:        token.LastSync,
				}); err != nil {
					util.Log.Println("error completing message sync for ", token.ID, ":", err)
				}
			}
		}
	}()

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"info":    convInfo,
		"missing": missingTokens,
	})
}
