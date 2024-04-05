package townsquare_handlers

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipeshandler"
)

// Action: townsquare_join
func joinTownsquare(message pipeshandler.Context) {

	// Get account from backend
	res, err := integration.PostRequest(integration.BasePath+"/v1/account/get_node", map[string]interface{}{
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

	caching.JoinTownsquare(message.Client.ID, res["name"].(string))
	pipeshandler.SuccessResponse(message)
}

// Action: townsquare_leave
func leaveTownsquare(message pipeshandler.Context) {
	// TODO: Implement
}
