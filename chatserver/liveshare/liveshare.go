package liveshare

import (
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/Liphium/station/chatserver/util"
)

type Transaction struct {
	Id               string
	UploadToken      string // Required to upload the file
	Token            string // Required to join the transaction
	VolumePath       string
	Account          string
	FileName         string
	PriorityReceiver string
	FileSize         int64
	CurrentIndex     int64
	Range            SendRange
	ReceiversCache   sync.Map
}

func (t *Transaction) ChunkFilePath(chunk int64) string {
	return fmt.Sprintf("%s/%s", t.VolumePath, t.ChunkFileName(chunk))
}

func (t *Transaction) ChunkFileName(chunk int64) string {
	return fmt.Sprintf("chunk_%d", chunk)
}

type TransactionReceiver struct {
	Mutex         *sync.Mutex
	TransactionId string
	ReceiverId    string
	CurrentIndex  int64
	Sent          bool
	SendChannel   chan *[]byte
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

	// Compute values
	endIndex := int64(math.Ceil(float64(fileSize) / float64(ChunkSize)))
	path := os.Getenv("CN_LS_REPO") + "/" + id

	transaction := &Transaction{
		Id:             id,
		UploadToken:    util.GenerateToken(50),
		Token:          util.GenerateToken(50),
		VolumePath:     path,
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
