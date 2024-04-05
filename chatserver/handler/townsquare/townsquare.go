package townsquare_handlers

import "github.com/Liphium/station/chatserver/caching"

func SetupActions() {
	caching.CSInstance.RegisterHandler("townsquare_join", joinTownsquare)
}
