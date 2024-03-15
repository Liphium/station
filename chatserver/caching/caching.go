package caching

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// This just needs to be kept somewhere to avoid import cycles
var Instance *pipeshandler.Instance
var Node *pipes.LocalNode

func SetupCaches() {
	setupConversationsCache()
	setupMembersCache()
	setupCallsCache()
	setupAdapterCache()
}
