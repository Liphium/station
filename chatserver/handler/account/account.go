package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

func SetupActions() {
	wshandler.RegisterHandler(caching.Node, "st_ch", changeStatus)
	wshandler.RegisterHandler(caching.Node, "st_send", sendStatus)
	wshandler.RegisterHandler(caching.Node, "st_res", respondToStatus)
}
