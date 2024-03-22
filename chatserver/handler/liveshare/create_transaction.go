package liveshare_actions

import (
	"github.com/Liphium/station/chatserver/liveshare"
	"github.com/Liphium/station/pipeshandler"
)

func createTransaction(context pipeshandler.Context) {

	if context.ValidateForm("name", "size") {
		pipeshandler.ErrorResponse(context, "invalid")
	}

	name := context.Data["name"].(string)
	size := int64(context.Data["size"].(float64))

	transaction, ok := liveshare.NewTransaction(context.Client.ID, name, size)
	if !ok {
		pipeshandler.ErrorResponse(context, "failed") // TODO: Better message
		return
	}

	pipeshandler.NormalResponse(context, map[string]interface{}{
		"success":      true,
		"id":           transaction.Id,
		"token":        transaction.Token,
		"upload_token": transaction.UploadToken,
	})
}
