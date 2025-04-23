package conversation_routes

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Public so it can be unit tested (in the future ig)
type OpenConversationRequest struct {
	AccountData string   `json:"accountData"` // Account data of the user opening the conversation (encrypted)
	Members     []string `json:"members"`
	Type        uint     `json:"type"`
	Data        string   `json:"data"` // Encrypted data
}

func (r *OpenConversationRequest) Validate() bool {
	return len(r.Data) > 0 && len(r.Data) <= util.MaxConversationDataLength
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
		return integration.FailedRequest(c, localization.ErrorGroupMemberLimit(util.MaxConversationMembers), nil)
	}

	if len(req.AccountData) > util.MaxConversationTokenDataLength {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	for _, member := range req.Members {
		if len(member) > util.MaxConversationTokenDataLength {
			return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
		}
	}

	// Check if the conversation type is valid
	if len(req.Members) > 1 && req.Type == database.ConvTypePrivateMessage {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}
	convType := database.ConvTypeGroup
	switch req.Type {
	case database.ConvTypePrivateMessage:
		convType = database.ConvTypePrivateMessage
	case database.ConvTypeGroup:
		convType = database.ConvTypeGroup
	case database.ConvTypeSquare:
		convType = database.ConvTypeSquare
	}

	// Generate the address for the conversation
	conv := database.Conversation{
		ID:      util.GenerateToken(util.ConversationIDLength) + "@" + integration.Domain,
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

		tk := database.ConversationToken{
			ID:           util.GenerateToken(util.ConversationTokenIDLength) + "@" + integration.Domain,
			Conversation: conv.ID,
			Activated:    false,
			Token:        convToken,
			Rank:         database.RankUser,
			Data:         memberData,
			Reads:        "",
		}

		if err := database.DBConn.Create(&tk).Error; err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		tokens[util.HashString(memberData)] = returnableToken{
			ID:    tk.ID,
			Token: convToken,
		}
	}

	adminToken := database.ConversationToken{
		ID:           util.GenerateToken(util.ConversationTokenIDLength) + "@" + integration.Domain,
		Token:        util.GenerateToken(util.ConversationTokenLength),
		Activated:    true,
		Conversation: conv.ID,
		Rank:         database.RankAdmin,
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
