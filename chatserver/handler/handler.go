package handler

import (
	"github.com/Liphium/station/chatserver/handler/account"
	"github.com/Liphium/station/chatserver/handler/conversation"
	zapshare_actions "github.com/Liphium/station/chatserver/handler/zapshare"
)

func Create() {
	conversation.SetupActions()
	account.SetupActions()
	zapshare_actions.SetupActions()
}
