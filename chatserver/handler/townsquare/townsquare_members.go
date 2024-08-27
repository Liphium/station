package townsquare_handlers

import (
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler"
)

// Action: townsquare_join
func joinTownsquare(message pipeshandler.Context) {

	// Get account from backend
	res, err := integration.PostRequestBackend("/account/get_node", map[string]interface{}{
		"id":    message.Client.ID,
		"node":  util.NodeTo64(caching.CSNode.ID),
		"token": caching.CSNode.Token,
	})

	if err != nil {
		pipeshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	if !res["success"].(bool) {
		util.Log.Println("something went wrong yk", res["error"].(string))
		pipeshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	// Join the square and init client
	caching.JoinTownsquare(message.Client.ID, res["name"].(string), res["sg"].(string))
	caching.SendAllTownsquareMembers(message.Client.ID)

	pipeshandler.SuccessResponse(message)
}

// Action: townsquare_leave
func leaveTownsquare(message pipeshandler.Context) {
	caching.LeaveTownsquare(message.Client.ID)
	pipeshandler.SuccessResponse(message)
}

// Action: townsquare_open
func openTownsquare(message pipeshandler.Context) {
	caching.SetTownsquareViewing(message.Client.ID, true)

	// Send messages over
	caching.SendMessages(message.Client.ID, time.Now().UnixMilli())

	pipeshandler.SuccessResponse(message)
}

// Action: townsquare_close
func closeTownsquare(message pipeshandler.Context) {
	caching.SetTownsquareViewing(message.Client.ID, false)
	pipeshandler.SuccessResponse(message)
}
