package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	pipeshandler.CreateHandlerFor(caching.CSInstance, "st_res", respondToStatus)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "st_validate", statusValidateAction)
}

// Do some basic status validation
func ValidateStatus(status string, data string) bool {
	if len(status) >= 1000 || len(data) >= 1000 {
		return false
	}
	return true
}

func StatusEvent(st string, data string, conversation string, ownToken string, suffix string) pipes.Event {
	return pipes.Event{
		Name: "acc_st" + suffix,
		Data: map[string]interface{}{
			"c":  conversation,
			"o":  ownToken,
			"st": st,
			"d":  data,
		},
	}
}
