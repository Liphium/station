package caching

import (
	"errors"
	"slices"
	"time"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/dgraph-io/ristretto"
)

// ! Always use cost 1
var conversationsCache *ristretto.Cache // Conversation token ID -> Conversation Token
const ConversationTTL = time.Hour * 1   // 1 hour

// Errors
var ErrInvalidToken = errors.New(localization.InvalidRequest)

func setupConversationsCache() {
	var err error

	// TODO: Check if values really are enough
	conversationsCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of objects expected (1,000,000).
		MaxCost:     1 << 30, // maximum cost of cache (1,000,000).
		BufferItems: 64,      // something from the docs
	})
	if err != nil {
		panic(err)
	}
}

// This does database requests and stuff
func ValidateToken(id string, token string) (conversations.ConversationToken, error) {

	// Check cache
	if value, found := conversationsCache.Get(id); found {

		// Check if token is valid
		if value.(conversations.ConversationToken).Token != token {
			return conversations.ConversationToken{}, ErrInvalidToken
		}

		return value.(conversations.ConversationToken), nil
	}

	var conversationToken conversations.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return conversations.ConversationToken{}, err
	}

	// Add to cache
	conversationsCache.SetWithTTL(id, conversationToken, 1, ConversationTTL)

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
		if value, found := conversationsCache.Get(token.ID); found {
			if value.(conversations.ConversationToken).Token == token.Token {
				foundTokens = append(foundTokens, value.(conversations.ConversationToken))
			}
			continue
		} else {
			notFound[token.ID] = token
			notFoundIds = append(notFoundIds, token.ID)
		}
	}

	// Get tokens from database
	var conversationTokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", notFoundIds).Find(&conversationTokens).Error; err != nil {
		return nil, nil, err
	}

	for _, token := range conversationTokens {
		conversationsCache.SetWithTTL(token.ID, token, 1, ConversationTTL)
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

// Does a lookup in the database on all tokens
func ValidateTokensLookup(tokens *[]conversations.SentConversationToken) ([]conversations.ConversationToken, []string, error) {

	tokenIds := make([]string, len(*tokens))
	for i, token := range *tokens {
		tokenIds[i] = token.ID
	}

	// Get tokens from database
	var conversationTokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("id IN ?", tokenIds).Find(&conversationTokens).Error; err != nil {
		return nil, nil, err
	}

	for _, token := range conversationTokens {
		conversationsCache.SetWithTTL(token.ID, token, 1, ConversationTTL)
	}

	// Get all missing tokens to delete those conversations from the client
	for _, token := range conversationTokens {
		index := slices.Index(tokenIds, token.ID)
		if index >= 0 {
			tokenIds = append(tokenIds[:index], tokenIds[index+1:]...)
		}
	}

	return conversationTokens, tokenIds, nil
}

// Deletes cache and does database queries again (for when caching would break something)
func GetTokenNew(id string) (conversations.ConversationToken, error) {
	conversationsCache.Del(id)
	return GetToken(id)
}

// Get a conversation token
func GetToken(id string) (conversations.ConversationToken, error) {

	// Check cache
	if value, found := conversationsCache.Get(id); found {
		return value.(conversations.ConversationToken), nil
	}

	var conversationToken conversations.ConversationToken
	if err := database.DBConn.Where("id = ?", id).Take(&conversationToken).Error; err != nil {
		return conversations.ConversationToken{}, err
	}

	// Add to cache
	conversationsCache.SetWithTTL(id, conversationToken, 1, ConversationTTL)

	return conversationToken, nil
}

func UpdateToken(token conversations.ConversationToken) error {

	// Update cache
	conversationsCache.SetWithTTL(token.ID, token, 1, ConversationTTL)

	return nil
}

func DeleteToken(id, token string) {
	CSNode.RemoveNodeWS("s-" + token)
	conversationsCache.Del(id)
}
