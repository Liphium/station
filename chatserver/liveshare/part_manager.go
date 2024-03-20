package liveshare

import (
	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/pipes"
)

// Register a new transaction receiver
func NewTransactionReceiver(id string, token string) (string, bool) {
	obj, ok := transactionsCache.Load(id)
	if !ok {
		return "", false
	}
	transaction := obj.(*Transaction)
	if transaction.Token != token {
		return "", false
	}

	receiverId := util.GenerateToken(10)
	for {
		_, ok := transaction.ReceiversCache.Load(receiverId)
		if !ok {
			break
		}
		receiverId = util.GenerateToken(10)
	}

	transaction.ReceiversCache.Store(receiverId, &TransactionReceiver{
		TransactionId: id,
		ReceiverId:    receiverId,
		CurrentIndex:  0,
		CurrentRange: SendRange{
			StartIndex: 0,
			EndIndex:   transaction.Range.EndIndex,
		},
		MissedRanges: []SendRange{},
	})

	return receiverId, true
}

// Send the uploader an update to upload more parts
func SendUploaderUpdate(id string) bool {

	transaction, ok := GetTransaction(id)
	if !ok {
		return false
	}

	// TODO: Compute the proper value here

	if err := caching.CSNode.SendClient(transaction.Account, pipes.Event{
		Name: "transaction_send_part",
		Data: map[string]interface{}{
			"index": transaction.CurrentIndex,
		},
	}); err != nil {
		return false
	}

	return true
}

func ReceiveNewPart() {
	// TODO
}
