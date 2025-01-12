package zapshare_actions

import (
	"os"

	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

type createTransactionAction struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func createTransaction(context *pipeshandler.Context, action createTransactionAction) pipes.Event {

	if os.Getenv("CHAT_NODE") == "" {
		util.Log.Println("Zap is disabled because CHAT_NODE is not set. It should be set to the URL of the chat node.")
		return pipeshandler.ErrorResponse(context, localization.ErrorInvalidRequest, nil)
	}

	// Create a new transaction
	transaction, ok := zapshare.NewTransaction(context.Client.ID, action.Name, action.Size)
	if !ok {
		return pipeshandler.ErrorResponse(context, localization.ErrorServer, nil)
	}

	return pipeshandler.NormalResponse(context, map[string]interface{}{
		"success":      true,
		"id":           transaction.Id,
		"token":        transaction.Token,
		"upload_token": transaction.UploadToken,
		"url":          os.Getenv("CHAT_NODE"),
	})
}
