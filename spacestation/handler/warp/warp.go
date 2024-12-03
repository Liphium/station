package warp_handlers

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandlers() {
	pipeshandler.CreateHandlerFor(caching.SSInstance, "wp_send", sendPacket)
}
