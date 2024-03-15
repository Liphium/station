package account

import "github.com/Liphium/station/chatserver/caching"

func SetupActions() {
	caching.Instance.RegisterHandler("st_ch", changeStatus)
	caching.Instance.RegisterHandler("st_send", sendStatus)
	caching.Instance.RegisterHandler("st_res", respondToStatus)
}
