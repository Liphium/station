package liveshare

import (
	"errors"
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

	receiverId := util.GenerateToken(10)
	for {
		_, ok := transaction.ReceiversCache.Load(receiverId)
		if !ok {
			break
		}
		receiverId = util.GenerateToken(10)
	}

	receiver := &TransactionReceiver{
		TransactionId: id,
		ReceiverId:    receiverId,
		CurrentIndex:  0,
		CurrentRange: SendRange{
			StartIndex: 0,
			EndIndex:   transaction.Range.EndIndex,
		},
		MissedRanges: []SendRange{},
		Mutex:        &sync.Mutex{},
	}
	transaction.ReceiversCache.Store(receiverId, receiver)

	return receiver, true
}

// Send the uploader an update to upload more parts
func (t *Transaction) RequestUploaderParts() bool {

	if err := caching.CSNode.SendClient(t.Account, pipes.Event{
		Name: "transaction_send_part",
		Data: map[string]interface{}{
			"index": t.CurrentIndex + ChunksAhead,
		},
	}); err != nil {
		return false
	}

	return true
}

// Called when a new part is received by any receiver
func (t *Transaction) PartReceived(id string, receiverId string) error {

	// TODO: Compute if the part can actually be deleted

	// Delete the part
	if err := os.Remove(t.ChunkFileName(t.CurrentIndex)); err != nil {
		return err
	}

	obj, ok := t.ReceiversCache.Load(receiverId)
	if !ok {
		return errors.New("receiver not found")
	}
	receiver := obj.(*TransactionReceiver)

	receiver.Mutex.Lock()
	defer receiver.Mutex.Unlock()

	if t.PriorityReceiver == receiver.ReceiverId {
		t.CurrentIndex++
		receiver.CurrentIndex++
		receiver.Sent = false
		t.RequestUploaderParts()
	}

	// Send a new part to the receiver
	file, err := os.Open(t.ChunkFilePath(t.CurrentIndex))
	if err != nil {
		// The next part will automatically be sent when the uploader sends the next part
		return nil
	}

	// Send the part
	part := make([]byte, ChunkSize)
	_, err = file.Read(part)
	if err != nil {
		return err
	}

	receiver.SendChannel <- &part
	receiver.Sent = true
	return nil
}

// Called when the uploader has uploaded a part
func (t *Transaction) PartUploaded(id string, index int64) error {

	t.ReceiversCache.Range(func(key, value any) bool {
		receiver := value.(*TransactionReceiver)
		if receiver.CurrentIndex == index && !receiver.Sent {
			receiver.SendChannel <- nil
		}
		return true
	})

	return nil
}

func (t *TransactionReceiver) SendPart(path string, index int64) {

}
