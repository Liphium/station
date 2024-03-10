package account

import (
	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/database/fetching"
	"github.com/Liphium/station/chatserver/util/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/send"
	"github.com/Liphium/station/pipeshandler/wshandler"
)

// Action: st_ch
func changeStatus(message wshandler.Message) {

	if !message.ValidateForm("status") {
		return
	}
	status := message.Data["status"].(string)

	// Save in database
	if err := database.DBConn.Model(&fetching.Status{}).Where("id = ?", message.Client.ID).Update("data", status).Error; err != nil {
		wshandler.ErrorResponse(message, localization.ErrorServer)
		return
	}

	// Send to all clients
	send.Client(message.Client.ID, pipes.Event{
		Name: "o:acc_st", // o: for own
		Data: map[string]interface{}{
			"st": status,
		},
	})

	wshandler.SuccessResponse(message)
}
