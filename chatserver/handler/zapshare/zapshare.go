package zapshare_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	caching.CSInstance.RegisterHandler("cancel_transaction", cancelTransaction)
	caching.CSInstance.RegisterHandler("create_transaction", createTransaction)
}

func cancelTransaction(ctx pipeshandler.Context) {
	zapshare.CancelTransactionByAccount(ctx.Client.ID)
	pipeshandler.SuccessResponse(ctx)
}
