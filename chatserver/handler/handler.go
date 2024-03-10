package handler

import (
	"github.com/Liphium/station/chatserver/handler/account"
	"github.com/Liphium/station/chatserver/handler/conversation"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

func Create() {
	wshandler.Initialize()

	conversation.SetupActions()
	account.SetupActions()
}
