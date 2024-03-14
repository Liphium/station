package space

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

func SetupActions() {
	wshandler.RegisterHandler(caching.Node, "spc_start", start)
	wshandler.RegisterHandler(caching.Node, "spc_join", joinCall)
	wshandler.RegisterHandler(caching.Node, "spc_leave", leaveCall)
}
