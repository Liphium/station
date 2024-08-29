package caching

import (
	"errors"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
)

// Errors
var ErrInvalidToken = errors.New(localization.InvalidRequest)

// This does database requests and stuff
func ValidateToken(id string, token string) (conversations.ConversationToken, error) {

	var conversationToken conversations.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return conversations.ConversationToken{}, err
	}

	if conversationToken.Token != token {
		return conversations.ConversationToken{}, ErrInvalidToken
	}

	return conversationToken, nil
}

// Returns: conversationTokens, missingTokens, tokenIds, err
func ValidateTokens(tokens *[]conversations.SentConversationToken) ([]conversations.ConversationToken, []string, []string, error) {
	foundTokens := []conversations.ConversationToken{}

	// Convert all tokens to data types they can be used with
	tokensMap := map[string]conversations.SentConversationToken{}
	tokenIds := []string{}
	for _, token := range *tokens {
		tokensMap[token.ID] = token
		tokenIds = append(tokenIds, token.ID)
	}

	// Get tokens from database
	var conversationTokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", tokenIds).Find(&conversationTokens).Error; err != nil {
		return nil, nil, nil, err
	}

	// Check if the tokens are actually there and sort them into the lists accordinly
	missingTokens := []string{}
	for _, token := range conversationTokens {
		if token.Token == tokensMap[token.ID].Token {

			// Add the token to the found tokens list
			tokensMap[token.ID] = conversations.SentConversationToken{
				ID:    "-",
				Token: "-",
			}
			foundTokens = append(foundTokens, token)
		} else {

			// Add the token to the missing tokens list
			missingTokens = append(missingTokens, token.ID)
			util.Log.Println("not found")
		}
	}

	return foundTokens, missingTokens, tokenIds, nil
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
