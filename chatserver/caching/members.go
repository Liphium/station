package caching

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
)

// TODO: Reimplement caching, but properly this time

type StoredMember struct {
	TokenID string // Conversation token ID
	Token   string // Conversation token
	Node    int64
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
				Node:    member.Node,
			}
		} else {
			storedMembers[i] = StoredMember{
				TokenID: member.ID,
				Token:   member.Token,
				Node:    member.Node,
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
				Node:    token.Node,
			})
		} else {
			returnMap[token.Conversation] = append(returnMap[token.Conversation], StoredMember{
				TokenID: token.ID,
				Token:   token.Token,
				Node:    token.Node,
			})
		}
	}

	return returnMap, nil
}

func MembersToPipes(members []StoredMember) ([]string, []string) {

	memberAdapters := make([]string, len(members))
	memberNodes := make([]string, len(members))

	for i, member := range members {
		memberAdapters[i] = "s-" + member.Token
		memberNodes[i] = util.Node64(member.Node)
	}

	return memberAdapters, memberNodes
}
