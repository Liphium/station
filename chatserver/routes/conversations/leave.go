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

type leaveRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Routes: /conversations/leave
func leaveConversation(c *fiber.Ctx) error {

	var req leaveRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	token, err := caching.ValidateToken(req.ID, req.Token)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("invalid conversation token: %s", err.Error()))
	}

	// Delete token
	if err := database.DBConn.Where("id = ?", token.ID).Delete(&conversations.ConversationToken{}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	caching.DeleteToken(token.ID, token.Token)

	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check if the chat is a DM (send delete message if it is)
	var conversation conversations.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Take(&conversation).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if conversation.Type == conversations.TypePrivateMessage && len(members) == 1 {

		// Delete the conversation
		if err := deleteConversation(conversation.ID); err != nil {
			return integration.FailedRequest(c, integration.ErrorServer, err)
		}

		return integration.SuccessfulRequest(c)
	}

	if len(members) == 0 {

		// Delete conversation
		if err := deleteConversation(conversation.ID); err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		return integration.SuccessfulRequest(c)
	} else {

		// Check if another admin is needed
		if token.Rank == conversations.RankAdmin {
			needed := true
			bestCase := conversations.ConversationToken{
				Rank: conversations.RankUser,
			}
			for _, member := range members {
				userToken, err := caching.GetToken(member.TokenID)
				if err != nil {
					continue
				}

				if userToken.Rank == conversations.RankAdmin {
					needed = false
					break
				}

				if bestCase.Rank <= userToken.Rank {
					bestCase = userToken
				}
			}

			// Promote to admin if needed
			if needed {
				if database.DBConn.Model(&conversations.ConversationToken{}).Where("id = ?", bestCase.ID).Update("rank", conversations.RankAdmin).Error != nil {
					return integration.FailedRequest(c, localization.ErrorServer, nil)
				}

				err = message_routes.SendSystemMessage(token.Conversation, message_routes.GroupNewAdmin, []string{message_routes.AttachAccount(bestCase.Data)})
				if err != nil {
					return integration.FailedRequest(c, localization.ErrorServer, nil)
				}
			}
		}
	}

	message_routes.SendSystemMessage(token.Conversation, message_routes.GroupMemberLeave, []string{
		message_routes.AttachAccount(token.Data),
	})

	return integration.SuccessfulRequest(c)
}
