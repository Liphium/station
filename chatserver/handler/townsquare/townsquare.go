package townsquare_handlers

import "github.com/Liphium/station/chatserver/caching"

func SetupActions() {
	caching.CSInstance.RegisterHandler("townsquare_join", joinTownsquare)
	caching.CSInstance.RegisterHandler("townsquare_leave", leaveTownsquare)
	caching.CSInstance.RegisterHandler("townsquare_open", openTownsquare)
	caching.CSInstance.RegisterHandler("townsquare_close", closeTownsquare)
	caching.CSInstance.RegisterHandler("townsquare_send", sendMessage)
	caching.CSInstance.RegisterHandler("townsquare_delete", deleteMessage)
}
