package caching

import (
	"errors"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
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

func ValidateTokens(tokens *[]conversations.SentConversationToken) ([]conversations.ConversationToken, []string, error) {

	// Check cache
	foundTokens := []conversations.ConversationToken{}

	notFound := map[string]conversations.SentConversationToken{}
	notFoundIds := []string{}
	for _, token := range *tokens {
		notFound[token.ID] = token
		notFoundIds = append(notFoundIds, token.ID)
	}

	// Get tokens from database
	var conversationTokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", notFoundIds).Find(&conversationTokens).Error; err != nil {
		return nil, nil, err
	}

	for _, token := range conversationTokens {
		if token.Token == notFound[token.ID].Token {
			notFound[token.ID] = conversations.SentConversationToken{
				ID:    "-",
				Token: "-",
			}
			foundTokens = append(foundTokens, token)
		}
	}

	// Get all missing tokens to delete those conversations from the client
	missingTokens := []string{}
	for id := range notFound {
		if id != "-" {
			missingTokens = append(missingTokens, id)
		}
	}

	return foundTokens, missingTokens, nil
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
