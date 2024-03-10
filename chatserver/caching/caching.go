package caching

func SetupCaches() {
	setupConversationsCache()
	setupMembersCache()
	setupCallsCache()
	setupAdapterCache()
}
