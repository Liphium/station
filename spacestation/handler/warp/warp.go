package warp_handlers

import (
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
)

func SetupHandlers() {
	pipeshandler.CreateHandlerFor(caching.SSInstance, "wp_send_to", sendPacketTo)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "wp_send_back", sendPacketBack)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "wp_disconnect", disconnect)
	pipeshandler.CreateHandlerFor(caching.SSInstance, "wp_create", create)
}
