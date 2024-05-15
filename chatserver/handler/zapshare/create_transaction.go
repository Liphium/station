package zapshare_actions

import (
	"os"

	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/pipeshandler"
)

func createTransaction(context pipeshandler.Context) {

	if os.Getenv("CHAT_NODE") == "" {
		util.Log.Println("Live share is disabled because CHAT_NODE is not set. It should be set to the URL of the chat node.")
		pipeshandler.ErrorResponse(context, "invalid")
		return
	}

	if context.ValidateForm("name", "size") {
		pipeshandler.ErrorResponse(context, "invalid")
	}

	name := context.Data["name"].(string)
	size := int64(context.Data["size"].(float64))

	transaction, ok := zapshare.NewTransaction(context.Client.ID, name, size)
	if !ok {
		pipeshandler.ErrorResponse(context, "failed") // TODO: Better message
		return
	}

	pipeshandler.NormalResponse(context, map[string]interface{}{
		"success":      true,
		"id":           transaction.Id,
		"token":        transaction.Token,
		"upload_token": transaction.UploadToken,
		"url":          os.Getenv("CHAT_NODE"),
	})
}
