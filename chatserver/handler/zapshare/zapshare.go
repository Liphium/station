package zapshare_actions

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/zapshare"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

func SetupActions() {
	pipeshandler.CreateHandlerFor(caching.CSInstance, "cancel_transaction", cancelTransaction)
	pipeshandler.CreateHandlerFor(caching.CSInstance, "create_transaction", createTransaction)
}

func cancelTransaction(ctx *pipeshandler.Context, data interface{}) pipes.Event {
	zapshare.CancelTransactionByAccount(ctx.Client.ID)
	return pipeshandler.SuccessResponse(ctx)
}
