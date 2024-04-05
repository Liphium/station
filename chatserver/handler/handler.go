package handler

import (
	"github.com/Liphium/station/chatserver/handler/account"
	"github.com/Liphium/station/chatserver/handler/conversation"
	liveshare_actions "github.com/Liphium/station/chatserver/handler/liveshare"
	townsquare_handlers "github.com/Liphium/station/chatserver/handler/townsquare"
)

func Create() {
	conversation.SetupActions()
	account.SetupActions()
	liveshare_actions.SetupActions()
	townsquare_handlers.SetupActions()
}
