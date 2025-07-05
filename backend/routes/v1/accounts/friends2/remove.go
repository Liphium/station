package friends2_routes

import (
	"fmt"

	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

/*
	// remove
		remove func
		send remove event to account
		send deletion request to remove_external

	// remove_external
		if external:
			check if accountLPH is from the server that sends the request
			remove func
			send remove event to account on server

	// remove func:
		- find friendship with ("target = ?", accountLPH)
		- find friendship with ("target = ?", targetLPH)
		- remove both if exist
		return error if 0 can be removed

*/

type FriendRemoveRequest struct {
	Id standards.LPHAddress `json:"id"` // Address of friend to remove
}

// Route: /a/accounts/friends/remove
func removeFriend(c *fiber.Ctx) error {
	var req FriendRemoveRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "request is invalid")
	}

	target := standards.LPHAddress(req.Id)
	_, targetTown, ok := target.Split()
	if !ok {
		return integration.FailedRequest(c, localization.ErrorInvalidRequestContent, nil)
	}

	sender := standards.LiphiumAddress(verify.InfoLocals(c).GetAccount())
	if err := deleteFriendAndRequest(sender, req.Id); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Send removal event to sender
	eventForSender := FriendRemoveEvent(target)

	// TODO: Register the adapter for account (in case decentralized) idk stand so in add external

	if err := service.Instance.SendOne(sender.String(), eventForSender); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, fmt.Errorf("couldn't send to sender: %s", err))
	}

	// If on a remote server, send the request to that server.
	if targetTown != standards.CurrentTown() {
		requests.PostRequest(targetTown, "/accounts/friends/remove_external", requests.Map{
			"from": standards.LiphiumAddress(verify.InfoLocals(c).GetAccount()),
			"to":   req.Id,
		})
	} else {
		// Send removal event to target
		eventForTarget := FriendRemoveEvent(sender)

		// TODO: Register the adapter for account (in case decentralized) idk stand so in add external

		if err := service.Instance.SendOne(target.String(), eventForTarget); err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, fmt.Errorf("couldn't send to target: %s", err))
		}
	}

	return integration.SuccessfulRequest(c)
}
