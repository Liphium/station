package zapshare

import (
	"errors"
	"math"
	"os"
	"sync"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/pipes"
)

// Register a new transaction receiver
func NewTransactionReceiver(id string, token string) (*TransactionReceiver, bool) {
	obj, ok := transactionsCache.Load(id)
	if !ok {
		return nil, false
	}
	transaction := obj.(*Transaction)
	if transaction.Token != token {
		return nil, false
	}

	if transaction.PriorityReceiver != "" {
		return nil, false
	}

	receiverId := util.GenerateToken(10)
	for {
		_, ok := transaction.ReceiversCache.Load(receiverId)
		if !ok {
			break
		}
		receiverId = util.GenerateToken(10)
	}

	if transaction.PriorityReceiver == "" {
		transaction.PriorityReceiver = receiverId
	}

	receiver := &TransactionReceiver{
		TransactionId: id,
		ReceiverId:    receiverId,
		CurrentIndex:  1,
		CurrentRange: SendRange{
			StartIndex: 1,
			EndIndex:   transaction.Range.EndIndex,
		},
		SendChannel:  make(chan int64),
		MissedRanges: []SendRange{},
		Mutex:        &sync.Mutex{},
		Waiting:      true,
	}
	transaction.ReceiversCache.Store(receiverId, receiver)

	transaction.RequestUploaderParts()
	return receiver, true
}

// Send the uploader an update to upload more parts
func (t *Transaction) RequestUploaderParts() bool {

	// Send the pipes event
	chunkStart := math.Min(float64(t.CurrentIndex), float64(t.Range.EndIndex))
	chunkEnd := math.Min(float64(t.CurrentIndex+ChunksAhead), float64(t.Range.EndIndex))
	if err := caching.CSNode.SendClient(t.Account, pipes.Event{
		Name: "transaction_send_part",
		Data: map[string]interface{}{
			"start": chunkStart,
			"end":   chunkEnd,
		},
	}); err != nil {
		return false
	}

	return true
}

// Called when a new part is received by any receiver (bool is if the transaction is finished)
func (t *Transaction) PartReceived(receiverId string) (bool, error) {

	// TODO: Compute if the part can actually be deleted

	// Delete the part
	if err := os.Remove(t.ChunkFilePath(t.CurrentIndex)); err != nil {
		return false, err
	}

	obj, ok := t.ReceiversCache.Load(receiverId)
	if !ok {
		return false, errors.New("receiver not found")
	}
	receiver := obj.(*TransactionReceiver)

	receiver.Mutex.Lock()
	defer receiver.Mutex.Unlock()

	if t.PriorityReceiver == receiver.ReceiverId {
		util.Log.Println("set as waiting")
		if t.CurrentIndex == t.Range.EndIndex {
			t.PriorityReceiver = ""
		}
		t.CurrentIndex++
		t.RequestUploaderParts()
	}

	if receiver.CurrentIndex == receiver.CurrentRange.EndIndex {
		receiver.SendChannel <- -1
		t.ReceiversCache.Delete(receiverId)
		return true, nil
	}
	receiver.CurrentIndex++

	return false, nil
}

// Called when the uploader has uploaded a part
func (t *Transaction) PartUploaded(index int64) error {

	t.ReceiversCache.Range(func(key, value any) bool {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					util.Log.Println(err)
				}
			}()

			receiver := value.(*TransactionReceiver)
			receiver.Mutex.Lock()
			defer receiver.Mutex.Unlock()

			receiver.SendChannel <- index
		}()
		return true
	})

	return nil
}
