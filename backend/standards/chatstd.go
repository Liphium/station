package standards

import (
	"regexp"
	"strings"
	"unsafe"
)

// Conversations
const ConversationIDLength = 12
const MaxConversationMembers = 100
const MaxConversationDataLength = 10_000 // Not sure about this number yet, but should be good for now
const MaxMessageExtra = 100

// Conversation tokens
const ConversationTokenIDLength = 12
const ConversationTokenLength = 32
const MaxConversationTokenDataLength = 2_000 // Not sure about this number yet, but should be good for now

func CheckMessageSize(message string) bool {
	return unsafe.Sizeof(message) > 1000*6
}

// Add the extra part added to the conversation id in the messages table (for topics in squares for example)
func MessageWithExtra(conversationId string, extra string) string {
	if extra == "" {
		return conversationId
	}
	return conversationId + "_" + extra
}

// Convert a message conversation id back to conversation id (first in return tuple) and extra
func IntoConversationAndExtra(conversation string) (string, string) {
	args := strings.Split(conversation, "_")
	if len(args) == 1 {
		return args[0], ""
	}
	return args[0], args[1]
}

// Regex for making sure there are only a-z, A-Z and 1-9 present in the extra part of the conversation id
const extraRegex = "^[a-zA-Z1-9]+$"

// Make sure the extra part isn't weird
func ValidateExtra(extra string) bool {
	if extra == "" {
		return true
	}
	if len(extra) > MaxMessageExtra {
		return false
	}
	matched, err := regexp.MatchString(extraRegex, extra)
	return err == nil && matched
}
