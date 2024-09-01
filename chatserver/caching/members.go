package caching

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/pipes"
)

// TODO: Reimplement caching, but properly this time

type StoredMember struct {
	TokenID string // Conversation token ID
	Token   string // Conversation token
}

// Does database requests and stuff
func LoadMembers(id string) ([]StoredMember, error) {

	var members []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation = ?", id).Find(&members).Error; err != nil {
		return []StoredMember{}, err
	}

	storedMembers := make([]StoredMember, len(members))
	for i, member := range members {
		if !member.Activated {
			storedMembers[i] = StoredMember{
				TokenID: member.ID,
				Token:   "-",
			}
		} else {
			storedMembers[i] = StoredMember{
				TokenID: member.ID,
				Token:   member.Token,
			}
		}
	}

	return storedMembers, nil
}

func LoadMembersArray(ids []string) (map[string][]StoredMember, error) {

	// Check cache
	returnMap := make(map[string][]StoredMember, len(ids)) // Conversation ID -> Members

	var tokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation IN ?", ids).Find(&tokens).Error; err != nil {
		return nil, err
	}
	for _, token := range tokens {
		if !token.Activated {
			returnMap[token.Conversation] = append(returnMap[token.Conversation], StoredMember{
				TokenID: token.ID,
				Token:   "-",
			})
		} else {
			returnMap[token.Conversation] = append(returnMap[token.Conversation], StoredMember{
				TokenID: token.ID,
				Token:   token.Token,
			})
		}
	}

	return returnMap, nil
}

// Send an event to all members in a conversation
func SendEventToMembers(members []StoredMember, event pipes.Event) error {

	// Make slices for a pipes send call
	memberAdapters := []string{}
	memberNodes := []string{}

	for i, member := range members {
		memberAdapters[i] = "s-" + member.TokenID
		memberNodes[i] = CSNode.ID
	}

	// Send event using pipes
	return CSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(memberAdapters, memberNodes),
		Event:   event,
	})
}
