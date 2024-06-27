package space

import (
	"github.com/Liphium/station/chatserver/caching"
)

func SetupActions() {
	caching.CSInstance.RegisterHandler("spc_start", start)
	caching.CSInstance.RegisterHandler("spc_join", joinCall)
	caching.CSInstance.RegisterHandler("spc_leave", leaveCall)
}
