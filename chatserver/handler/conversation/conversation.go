package conversation

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/handler/conversation/space"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	space.SetupActions()

	pipeshandler.CreateHandlerFor(caching.CSInstance, "conv_sub", subscribe)
}
