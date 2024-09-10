package conversation_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	message_actions "github.com/Liphium/station/chatserver/routes/actions/messages"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_leave
func HandleLeave(c *fiber.Ctx, token conversations.ConversationToken, _ interface{}) error {

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
		if err := action_helpers.DeleteConversation(conversation.ID); err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		return integration.SuccessfulRequest(c)
	}

	if len(members) == 0 {

		// Delete conversation
		if err := action_helpers.DeleteConversation(conversation.ID); err != nil {
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

				err = message_actions.SendSystemMessage(token.Conversation, message_actions.GroupNewAdmin, []string{message_actions.AttachAccount(bestCase.Data)})
				if err != nil {
					return integration.FailedRequest(c, localization.ErrorServer, nil)
				}
			}
		}
	}

	message_actions.SendSystemMessage(token.Conversation, message_actions.GroupMemberLeave, []string{
		message_actions.AttachAccount(token.Data),
	})

	return integration.SuccessfulRequest(c)
}
