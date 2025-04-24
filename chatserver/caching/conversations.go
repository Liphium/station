package caching

import (
	"errors"
	"slices"

	"github.com/Liphium/station/chatserver/database"
)

// This does database requests and stuff
func ValidateToken(id string, token string) (database.ConversationToken, error) {

	var conversationToken database.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return database.ConversationToken{}, err
	}

	if conversationToken.Token != token {
		return database.ConversationToken{}, errors.New("token is invalid")
	}

	return conversationToken, nil
}

// Returns: conversationTokens, missingTokens, conversationIds, err
func ValidateTokens(tokens *[]database.SentConversationToken) ([]database.ConversationToken, []string, []string, error) {
	foundTokens := []database.ConversationToken{}

	// Convert all tokens to data types they can be used with
	tokenIds := make([]string, len(*tokens))
	tokensMap := map[string]database.SentConversationToken{}
	for i, token := range *tokens {
		tokensMap[token.ID] = token
		tokenIds[i] = token.ID
	}

	// Get tokens from database
	var conversationTokens []database.ConversationToken
	if err := database.DBConn.Model(&database.ConversationToken{}).Where("id IN ?", tokenIds).Find(&conversationTokens).Error; err != nil {
		return nil, nil, nil, err
	}

	// Check if the tokens are actually there and sort them into the lists accordinly
	missingTokens := tokenIds
	conversationIds := []string{}
	for _, token := range conversationTokens {
		if token.Token == tokensMap[token.ID].Token {

			// Add the token to the found tokens list
			token.LastSync = tokensMap[token.ID].LastMessage
			tokensMap[token.ID] = database.SentConversationToken{
				ID:           "-",
				Token:        "-",
				Conversation: "-",
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
func GetToken(id string) (database.ConversationToken, error) {

	var conversationToken database.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return database.ConversationToken{}, err
	}

	return conversationToken, nil
}

func DeleteToken(id, token string) {
	CSNode.RemoveNodeWS("s-" + token)
}
