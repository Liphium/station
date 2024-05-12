package handler

import (
	"github.com/Liphium/station/chatserver/handler/account"
	"github.com/Liphium/station/chatserver/handler/conversation"
	townsquare_handlers "github.com/Liphium/station/chatserver/handler/townsquare"
	zapshare_actions "github.com/Liphium/station/chatserver/handler/zapshare"
)

func Create() {
	conversation.SetupActions()
	account.SetupActions()
	zapshare_actions.SetupActions()
	townsquare_handlers.SetupActions()
}
