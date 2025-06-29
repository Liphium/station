package friends2_routes

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Route: /accounts/friends/add
func addFriend(c *fiber.Ctx) error {
	var req struct {
		Id         string `json:"id"`
		ProfileKey string `json:"prf"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "request is invalid")
	}

	targetId, targetTown, ok := standards.LiphiumAddress(req.Id).Split()
	if !ok {
		return integration.FailedRequest(c, localization.ErrorInvalidRequestContent, nil)
	}

	// If on a remote server, send the request to that server. Otherwise handle here.
	if targetTown != standards.CurrentTown() {
		res, err := requests.PostRequest(targetTown, "/accounts/friends/add_external", requests.Map{
			"from": standards.LiphiumAddress(verify.InfoLocals(c).GetAccount()),
			"prf":  req.ProfileKey,
			"to":   req.Id,
		})
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorOtherServer, err)
		}
		if !requests.ValueOr(res, "success", false) {
			return integration.FailedRequest(c, localization.ErrorFromOtherServer(requests.ValueOr(res, "message", "?")), nil)
		}

		// Mirror the database entry in case a new one was created
		if requests.ValueOr(res, "created", false) {
			if err := database.DBConn.Create(&database.Friendship{
				Token:      requests.ValueOr(res, "token", ""),
				Request:    false,
				Account:    standards.LiphiumAddress(verify.InfoLocals(c).GetAccount()).String(),
				Target:     req.Id,
				ProfileKey: req.ProfileKey,
			}).Error; err != nil {
				// TODO: Call remove on the other server (spec not completed, so waiting for final version)
				return integration.FailedRequest(c, localization.ErrorServer, err)
			}
		}
	} else {
		target, err := uuid.Parse(targetId)
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorInvalidRequestContent, err)
		}

		address := standards.LiphiumAddress(verify.InfoLocals(c).GetAccount())
		if _, msg, err := createFriendRequest(address, target, req.ProfileKey); err != nil {
			if msg == nil {
				msg = localization.ErrorServer
			}
			return integration.FailedRequest(c, msg, err)
		}
	}

	return integration.SuccessfulRequest(c)
}
