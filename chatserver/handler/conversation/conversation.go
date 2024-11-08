package conversation

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	pipeshandler.CreateHandlerFor(caching.CSInstance, "conv_sub", subscribe)
}
