package zapshare

import (
	"log"
	"math"
	"sync"

	"github.com/Liphium/station/chatserver/caching"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/pipes"
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
	ReceiversCache   *sync.Map
	FileParts        *sync.Map // Chunk id -> file part
}

type TransactionReceiver struct {
	Mutex         *sync.Mutex
	TransactionId string
	ReceiverId    string
	CurrentIndex  int64
	Waiting       bool
	SendChannel   chan int64
	CurrentRange  SendRange
	MissedRanges  []SendRange
}

type SendRange struct {
	StartIndex int64
	EndIndex   int64
}

const ChunksAhead = 20
const ChunkSize = 1024 * 1024           // 1 MB
const MaxChunkSize = ChunkSize + 4*1024 // 1 MB + 4 KB (actual chunk is 1 MB, but there are additional headers for encryption)

// SessionId -> Transaction ID
var userTransactions sync.Map = sync.Map{}

// Transaction ID -> Transaction
var transactionsCache sync.Map = sync.Map{}

func NewTransaction(account string, fileName string, fileSize int64) (*Transaction, bool) {

	if userId, ok := userTransactions.Load(account); ok {
		CancelTransaction(userId.(string))
	}

	id := util.GenerateToken(10)
	for {
		_, ok := transactionsCache.Load(id)
		if !ok {
			break
		}
		id = util.GenerateToken(10)
	}

	// Compute values
	endIndex := int64(math.Ceil(float64(fileSize) / float64(ChunkSize)))
	log.Println("End index:", endIndex)

	transaction := &Transaction{
		Id:             id,
		UploadToken:    util.GenerateToken(50),
		Token:          util.GenerateToken(50),
		Account:        account,
		FileName:       fileName,
		FileSize:       fileSize,
		CurrentIndex:   1,
		Range:          SendRange{StartIndex: 1, EndIndex: endIndex},
		ReceiversCache: &sync.Map{},
		FileParts:      &sync.Map{},
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

func CancelTransactionByAccount(account string) {
	userId, ok := userTransactions.Load(account)
	if !ok {
		return
	}
	CancelTransaction(userId.(string))
}

func CancelTransaction(id string) {

	transaction, ok := GetTransaction(id)
	if !ok {
		return
	}

	// Disconnect all receivers
	transaction.ReceiversCache.Range(func(key, value interface{}) bool {
		util.Log.Println("Disconnecting receiver", key)
		receiver := value.(*TransactionReceiver)
		receiver.SendChannel <- -1
		return true
	})

	// Delete the transaction from the cache
	transactionsCache.Delete(id)
	userTransactions.Delete(transaction.Account)

	// Inform the sender
	caching.CSNode.SendClient(transaction.Account, pipes.Event{
		Name: "transaction_end",
		Data: map[string]interface{}{},
	})
}
