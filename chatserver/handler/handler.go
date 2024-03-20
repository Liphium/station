package handler

import (
	"github.com/Liphium/station/chatserver/handler/account"
	"github.com/Liphium/station/chatserver/handler/conversation"
	liveshare_actions "github.com/Liphium/station/chatserver/handler/liveshare"
)

func Create() {
	conversation.SetupActions()
	account.SetupActions()
	liveshare_actions.SetupActions()
}
