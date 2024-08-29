package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	pipeshandler.CreateHandlerFor(caching.CSInstance, "st_send", sendStatus)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "st_res", respondToStatus)
}
