package remote_event_channel

import (
	"github.com/Liphium/station/chatserver/caching"
	action_helpers "github.com/Liphium/station/chatserver/routes/actions/helpers"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/pipes"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

type remoteEventRequest struct {
	ID    string      `json:"id"`
	Token string      `json:"token"`
	Event pipes.Event `json:"event"`
}

// Route: /event_channel/send
func HandleRemoteEvent(c *fiber.Ctx) error {

	// Parse the request
	var req remoteEventRequest
	if err := integration.BodyParser(c, &req); err != nil {
		return integration.InvalidRequest(c, "request not valid")
	}

	// Check if the token was subscribed to from the current node
	obj, valid := action_helpers.TokenMap.Load(req.ID)
	if !valid {
		return integration.InvalidRequest(c, "token wasn't found")
	}
	data := obj.(*action_helpers.TokenData)
	if data.Token != req.Token {
		return integration.InvalidRequest(c, "token isn't valid")
	}

	// Marshal the event
	message, err := sonic.Marshal(req.Event)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send it to the adapter
	if err := caching.CSNode.AdapterReceiveWeb("s-"+req.ID, req.Event, message); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
