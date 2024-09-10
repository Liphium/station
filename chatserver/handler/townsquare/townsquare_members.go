package townsquare_handlers

import (
	"errors"
	"time"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: townsquare_join
func joinTownsquare(message *pipeshandler.Context, data interface{}) pipes.Event {

	// Get account from backend
	res, err := integration.PostRequestBackend("/account/get_node", map[string]interface{}{
		"id":    message.Client.ID,
		"node":  util.NodeTo64(caching.CSNode.ID),
		"token": caching.CSNode.Token,
	})

	if err != nil {
		return pipeshandler.ErrorResponse(message, localization.ErrorServer, err)
	}

	if !res["success"].(bool) {
		util.Log.Println("something went wrong yk", res["error"].(string))
		return pipeshandler.ErrorResponse(message, localization.ErrorServer, errors.New(res["error"].(string)))
	}

	// Join the square and init client
	caching.JoinTownsquare(message.Client.ID, res["name"].(string), res["sg"].(string))
	caching.SendAllTownsquareMembers(message.Client.ID)

	return pipeshandler.SuccessResponse(message)
}

// Action: townsquare_leave
func leaveTownsquare(message *pipeshandler.Context, data interface{}) pipes.Event {
	caching.LeaveTownsquare(message.Client.ID)
	return pipeshandler.SuccessResponse(message)
}

// Action: townsquare_open
func openTownsquare(message *pipeshandler.Context, data interface{}) pipes.Event {
	caching.SetTownsquareViewing(message.Client.ID, true)

	// Send messages over
	caching.SendMessages(message.Client.ID, time.Now().UnixMilli())

	return pipeshandler.SuccessResponse(message)
}

// Action: townsquare_close
func closeTownsquare(message *pipeshandler.Context, data interface{}) pipes.Event {
	caching.SetTownsquareViewing(message.Client.ID, false)
	return pipeshandler.SuccessResponse(message)
}
