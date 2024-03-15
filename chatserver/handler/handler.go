package handler

import (
	"github.com/Liphium/station/chatserver/handler/account"
	"github.com/Liphium/station/chatserver/handler/conversation"
)

func Create() {
	conversation.SetupActions()
	account.SetupActions()
}
