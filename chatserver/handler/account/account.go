package account

import "github.com/Liphium/station/chatserver/caching"

func SetupActions() {
	caching.CSInstance.RegisterHandler("st_ch", changeStatus)
	caching.CSInstance.RegisterHandler("st_send", sendStatus)
	caching.CSInstance.RegisterHandler("st_res", respondToStatus)
}
