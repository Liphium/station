package caching

import (
	"time"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/chatserver/util"
	"github.com/dgraph-io/ristretto"
)

// ! Always use cost 1
var membersCache *ristretto.Cache // Conversation ID -> Members
const MemberTTL = time.Hour * 1   // 1 hour

func setupMembersCache() {
	var err error

	// TODO: Check if values really are enough
	membersCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of objects expected (1,000,000).
		MaxCost:     1 << 30, // maximum cost of cache (1,000,000).
		BufferItems: 64,      // something from the docs
	})
	if err != nil {
		panic(err)
	}
}

type StoredMember struct {
	TokenID string // Conversation token ID
	Token   string // Conversation token
	Node    int64
}

// The value of a token in the cache if it should be relearned from the database
const actionRelearnToken = "reget"

// Always does database requests (use where caching would break stuff)
func LoadMembersNew(id string) ([]StoredMember, error) {
	membersCache.Del(id)
	return LoadMembers(id)
}

// Does database requests and stuff
func LoadMembers(id string) ([]StoredMember, error) {

	// Check cache
	if value, found := membersCache.Get(id); found {
		members := value.([]StoredMember)
		changes := false
		for i, member := range members {
			if member.Token == actionRelearnToken {
				var token string
				if err := database.DBConn.Model(&conversations.ConversationToken{}).Select("token").Where("conversation = ?", id).Take(&token).Error; err != nil {
					return []StoredMember{}, err
				}
				member.Token = token
				members[i] = member
				changes = true
			}
		}

		if changes {
			membersCache.SetWithTTL(id, members, 1, MemberTTL)
		}

		return members, nil
	}

	var members []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation = ?", id).Find(&members).Error; err != nil {
		return []StoredMember{}, err
	}

	storedMembers := make([]StoredMember, len(members))
	for i, member := range members {
		if !member.Activated {
			storedMembers[i] = StoredMember{
				TokenID: member.ID,
				Token:   actionRelearnToken,
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

	// Add to cache
	membersCache.SetWithTTL(id, storedMembers, 1, MemberTTL)

	return storedMembers, nil
}

func LoadMembersArray(ids []string) (map[string][]StoredMember, error) {

	// Check cache
	returnMap := make(map[string][]StoredMember, len(ids)) // Conversation ID -> Members
	notFound := []string{}

	for _, id := range ids {
		if value, found := membersCache.Get(id); found {
			returnMap[id] = value.([]StoredMember)
		} else {
			notFound = append(notFound, id)
		}
	}

	var tokens []conversations.ConversationToken
	if err := database.DBConn.Model(&conversations.ConversationToken{}).Where("conversation IN ?", notFound).Find(&tokens).Error; err != nil {
		return nil, err
	}
	for _, token := range tokens {
		if !token.Activated {
			returnMap[token.Conversation] = append(returnMap[token.Conversation], StoredMember{
				TokenID: token.ID,
				Token:   actionRelearnToken,
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
	for key, memberTokens := range returnMap {
		membersCache.SetWithTTL(key, memberTokens, 1, MemberTTL)
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
