package conversation

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/handler/conversation/space"
)

func SetupActions() {
	space.SetupActions()

	caching.CSInstance.RegisterHandler("conv_sub", subscribe)
}
