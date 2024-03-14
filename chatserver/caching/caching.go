package caching

import "github.com/Liphium/station/pipes"

// This just needs to be kept somewhere to avoid import cycles
var Node *pipes.LocalNode

func SetupCaches() {
	setupConversationsCache()
	setupMembersCache()
	setupCallsCache()
	setupAdapterCache()
}
