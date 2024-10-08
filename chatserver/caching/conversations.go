package caching

import (
	"errors"
	"slices"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
)

// This does database requests and stuff
func ValidateToken(id string, token string) (conversations.ConversationToken, error) {

	var conversationToken conversations.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return conversations.ConversationToken{}, err
	}

	if conversationToken.Token != token {
		return conversations.ConversationToken{}, errors.New("token is invalid")
	}

	return conversationToken, nil
}

// Returns: conversationTokens, missingTokens, conversationIds, err
func ValidateTokens(tokens *[]conversations.SentConversationToken) ([]conversations.ConversationToken, []string, []string, error) {
	foundTokens := []conversations.ConversationToken{}

	// Convert all tokens to data types they can be used with
	tokenIds := make([]string, len(*tokens))
	tokensMap := map[string]conversations.SentConversationToken{}
	for i, token := range *tokens {
		tokensMap[token.ID] = token
		tokenIds[i] = token.ID
	}

	// Get tokens from database
	var conversationTokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", tokenIds).Find(&conversationTokens).Error; err != nil {
		return nil, nil, nil, err
	}

	// Check if the tokens are actually there and sort them into the lists accordinly
	missingTokens := tokenIds
	conversationIds := []string{}
	for _, token := range conversationTokens {
		if token.Token == tokensMap[token.ID].Token {

			// Add the token to the found tokens list
			tokensMap[token.ID] = conversations.SentConversationToken{
				ID:    "-",
				Token: "-",
			}
			foundTokens = append(foundTokens, token)

			// Delete the token from the missing tokens slice to make sure it isn't deleted
			missingTokens = slices.DeleteFunc(missingTokens, func(element string) bool {
				return element == token.ID
			})
		}
		conversationIds = append(conversationIds, token.Conversation)
	}

	return foundTokens, missingTokens, conversationIds, nil
}

// Get a conversation token
func GetToken(id string) (conversations.ConversationToken, error) {

	var conversationToken conversations.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return conversations.ConversationToken{}, err
	}

	return conversationToken, nil
}

func DeleteToken(id, token string) {
	CSNode.RemoveNodeWS("s-" + token)
}
