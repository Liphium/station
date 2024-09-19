package townsquare_handlers

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	pipeshandler.CreateHandlerFor(caching.CSInstance, "townsquare_join", joinTownsquare)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "townsquare_leave", leaveTownsquare)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "townsquare_open", openTownsquare)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "townsquare_close", closeTownsquare)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "townsquare_send", sendMessage)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "townsquare_delete", deleteMessage)
}
