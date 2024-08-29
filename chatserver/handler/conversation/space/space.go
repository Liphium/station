package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	pipeshandler.CreateHandlerFor(caching.CSInstance, "spc_start", start)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "spc_join", joinCall)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "spc_leave", leaveCall)
}
