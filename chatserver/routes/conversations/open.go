package conversation_routes

import (
	"os"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/gofiber/fiber/v2"
)

// Public so it can be unit tested (in the future ig)
type OpenConversationRequest struct {
	AccountData string   `json:"accountData"` // Account data of the user opening the conversation (encrypted)
	Members     []string `json:"members"`
	Data        string   `json:"data"` // Encrypted data
}

func (r *OpenConversationRequest) Validate() bool {
	return len(r.Members) > 0 && len(r.Data) > 0 && len(r.Data) <= util.MaxConversationDataLength
}

type returnableToken struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// Route: /conversations/open
func openConversation(c *fiber.Ctx) error {

	var req OpenConversationRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, err.Error())
	}

	// Validate request
	if !req.Validate() {
		return integration.InvalidRequest(c, "request is invalid")
	}

	if len(req.Members)+1 > util.MaxConversationMembers {
		return integration.FailedRequest(c, "member.limit", nil)
	}

	if len(req.AccountData) > util.MaxConversationTokenDataLength {
		return integration.FailedRequest(c, "data.limit", nil)
	}

	for _, member := range req.Members {
		if len(member) > util.MaxConversationTokenDataLength {
			return integration.FailedRequest(c, "data.limit", nil)
		}
	}

	// Determine the conversation type
	convType := conversations.TypePrivateMessage
	if len(req.Members) > 1 {
		convType = conversations.TypeGroup
	}

	// Generate the address for the conversation
	conv := conversations.Conversation{
		ID:      util.GenerateToken(util.ConversationIDLength) + "@" + os.Getenv("PROTOCOL") + os.Getenv("BASE_PATH"),
		Type:    uint(convType),
		Version: 1,
		Data:    req.Data,
	}

	if err := database.DBConn.Create(&conv).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	// Create tokens
	var tokens map[string]returnableToken = make(map[string]returnableToken)
	for _, memberData := range req.Members {

		convToken := util.GenerateToken(util.ConversationTokenLength)

		tk := conversations.ConversationToken{
			ID:           util.GenerateToken(util.ConversationTokenIDLength),
			Conversation: conv.ID,
			Activated:    false,
			Token:        convToken,
			Rank:         conversations.RankUser,
			Data:         memberData,
			LastRead:     0,
		}

		if err := database.DBConn.Create(&tk).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		tokens[util.HashString(memberData)] = returnableToken{
			ID:    tk.ID,
			Token: convToken,
		}
	}

	adminToken := conversations.ConversationToken{
		ID:           util.GenerateToken(util.ConversationTokenIDLength),
		Token:        util.GenerateToken(util.ConversationTokenLength),
		Activated:    true,
		Conversation: conv.ID,
		Rank:         conversations.RankAdmin,
		Data:         req.AccountData,
	}

	if err := database.DBConn.Create(&adminToken).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success":      true,
		"conversation": conv.ID,
		"type":         conv.Type,
		"admin_token": returnableToken{
			ID:    adminToken.ID,
			Token: adminToken.Token,
		},
		"tokens": tokens,
	})
}
