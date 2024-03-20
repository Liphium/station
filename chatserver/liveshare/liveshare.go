package liveshare

import (
	"math"
	"sync"

	"github.com/Liphium/station/chatserver/util"
)

type Transaction struct {
	Id               string
	UploadToken      string // Required to upload the file
	Token            string // Required to join the transaction
	Account          string
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
	CurrentIndex  int64
	CurrentRange  SendRange
	MissedRanges  []SendRange
}

type SendRange struct {
	StartIndex int64
	EndIndex   int64
}

const ChunksAhead = 10
const ChunkSize = 512 * 1024 // 512KB

// SessionId -> Transaction ID
var userTransactions sync.Map = sync.Map{}

// Transaction ID -> Transaction
var transactionsCache sync.Map = sync.Map{}

func NewTransaction(account string, fileName string, fileSize int64) (*Transaction, bool) {

	if _, ok := userTransactions.Load(account); ok {
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
		UploadToken:    util.GenerateToken(50),
		Token:          util.GenerateToken(50),
		Account:        account,
		FileName:       fileName,
		FileSize:       fileSize,
		Range:          SendRange{StartIndex: 0, EndIndex: endIndex},
		ReceiversCache: sync.Map{},
	}
	transactionsCache.Store(id, transaction)
	userTransactions.Store(account, id)

	return transaction, true
}

func GetTransaction(id string) (*Transaction, bool) {
	obj, ok := transactionsCache.Load(id)
	if !ok {
		return nil, false
	}
	return obj.(*Transaction), true
}
