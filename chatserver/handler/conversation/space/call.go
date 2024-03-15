package space

import (
	"github.com/Liphium/station/chatserver/caching"
)

func SetupActions() {
	caching.Instance.RegisterHandler("spc_start", start)
	caching.Instance.RegisterHandler("spc_join", joinCall)
	caching.Instance.RegisterHandler("spc_leave", leaveCall)
}
