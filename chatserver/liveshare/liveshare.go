package liveshare

import (
	"math"
	"sync"

	"github.com/Liphium/station/chatserver/util"
)

type Transaction struct {
	Id               string
	Token            string // Required to join the transaction
	Session          string
	FileName         string
	PriorityReceiver string
	FileSize         int64
	CurrentIndex     int64
	Range            SendRange
	ReceiversCache   sync.Map
}

type TransactionReceiver struct {
	TransactionId string
	ReceiverId    string
	CurrentRange  SendRange
	MissedRanges  []SendRange
}

type SendRange struct {
	StartIndex int64
	EndIndex   int64
}

const ChunksAhead = 3
const ChunkSize = 512 * 1024 // 512KB

// SessionId -> Transaction ID
var userTransactions sync.Map = sync.Map{}

// Transaction ID -> Transaction
var transactionsCache sync.Map = sync.Map{}

func NewTransaction(session string, fileName string, fileSize int64) (*Transaction, bool) {

	if _, ok := userTransactions.Load(session); ok {
		return nil, false
	}

	id := util.GenerateToken(10)
	for {
		_, ok := transactionsCache.Load(id)
		if !ok {
			break
		}
		id = util.GenerateToken(10)
	}

	// Compute range
	endIndex := int64(math.Ceil(float64(fileSize) / float64(ChunkSize)))

	transaction := &Transaction{
		Id:             id,
		Session:        session,
		FileName:       fileName,
		FileSize:       fileSize,
		Range:          SendRange{StartIndex: 0, EndIndex: endIndex},
		ReceiversCache: sync.Map{},
	}
	transactionsCache.Store(id, transaction)
	userTransactions.Store(session, id)

	return transaction, true
}

func NewTransactionReceiver(id string, token string, receiverAdapter string) {
	obj, ok := transactionsCache.Load(id)
	if !ok {
		return
	}
	transaction := obj.(*Transaction)
	if transaction.Token == token {
		transaction.ReceiversCache.Store(receiverAdapter, &TransactionReceiver{
			TransactionId: id,
			ReceiverId:    receiverAdapter,
			CurrentRange: SendRange{
				StartIndex: 0,
				EndIndex:   transaction.Range.EndIndex,
			},
			MissedRanges: []SendRange{},
		})
	}
}
