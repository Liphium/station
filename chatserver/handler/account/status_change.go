package account

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/fetching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: st_ch
func changeStatus(ctx pipeshandler.Context) {

	if !ctx.ValidateForm("status") {
		return
	}
	status := ctx.Data["status"].(string)

	// Save in database
	if err := database.DBConn.Model(&fetching.Status{}).Where("id = ?", ctx.Client.ID).Update("data", status).Error; err != nil {
		pipeshandler.ErrorResponse(ctx, localization.ErrorServer)
		return
	}

	// Send to all clients
	caching.CSNode.SendClient(ctx.Client.ID, pipes.Event{
		Name: "o:acc_st", // o: for own
		Data: map[string]interface{}{
			"st": status,
		},
	})

	pipeshandler.SuccessResponse(ctx)
}
