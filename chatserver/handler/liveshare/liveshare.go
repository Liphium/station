package liveshare_actions

import "github.com/Liphium/station/chatserver/caching"

func SetupActions() {
	caching.CSInstance.RegisterHandler("create_transaction", createTransaction)
}
