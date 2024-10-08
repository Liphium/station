package account

import (
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: st_validate
func statusValidateAction(c *pipeshandler.Context, action struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}) pipes.Event {

	// Do some basic status validation
	if !ValidateStatus(action.Status, action.Data) {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	return pipeshandler.SuccessResponse(c)
}
