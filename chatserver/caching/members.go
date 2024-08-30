package caching

import (
	"fmt"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/conversations"
	"github.com/Liphium/station/pipes"
)

// TODO: Reimplement caching, but properly this time

type StoredMember struct {
	TokenID string // Conversation token ID
	Token   string // Conversation token
	Remote  bool   // Whether the guy is connected remote or not
	Node    string // The domain or id of the connected node
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

func SendEventToMembers(members []StoredMember, event pipes.Event) error {

	// Make slices for a pipes send call
	memberAdapters := []string{}
	memberNodes := []string{}

	// Make a slice for members that need to be contacted using a remote event
	remoteMembers := []StoredMember{}

	for i, member := range members {
		if !member.Remote {
			// Add them to the pipes send when they are not from a remote instance
			memberAdapters[i] = "s-" + member.Token
			memberNodes[i] = member.Node
		} else {
			// Let the event be sent remotely
			remoteMembers = append(remoteMembers, member)
		}
	}

	// Send event using pipes
	err := CSNode.Pipe(pipes.ProtocolWS, pipes.Message{
		Channel: pipes.Conversation(memberAdapters, memberNodes),
		Event:   event,
	})
	if err != nil {
		return err
	}

	// Send the event to all the people who are connected remotely
	for _, member := range remoteMembers {
		fmt.Printf("member.Remote: %v\n", member.Remote)
	}

	return nil
}
