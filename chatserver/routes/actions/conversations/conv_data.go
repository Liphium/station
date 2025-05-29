package conversation_actions

import (
	"fmt"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type returnableMemberToken struct {
	ID   string `json:"id"`   // Conversation token id
	Data string `json:"data"` // Account id (encrypted)
	Rank uint   `json:"rank"` // Conversation rank
}

// Action: conv_data
func HandleGetData(c *fiber.Ctx, token database.ConversationToken, _ interface{}) error {

	// Get the conversation from the database
	var conversation database.Conversation
	if err := database.DBConn.Where("id = ?", token.Conversation).Take(&conversation).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// We use methods without caching here because if a member leaves on a different node, the cache won't be cleared
	members, err := caching.LoadMembers(token.Conversation)
	if err != nil {
		return integration.InvalidRequest(c, fmt.Sprintf("couldn't load members: %s", err.Error()))
	}

	realMembers := make([]returnableMemberToken, len(members))
	for i, memberToken := range members {

		member, err := caching.GetToken(memberToken.TokenID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		realMembers[i] = returnableMemberToken{
			ID:   member.ID,
			Data: member.Data,
			Rank: member.Rank,
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"id":      conversation.ID,
		"version": conversation.Version,
		"data":    conversation.Data,
		"members": realMembers,
	})
}
