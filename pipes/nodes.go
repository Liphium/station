package pipes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"sync"

	"github.com/dgraph-io/ristretto"
)

type Node struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	WS    string `json:"ws,omitempty"` // Websocket ip
	SL    string `json:"sl,omitempty"` // Socketless pipe

	// Encryption
	Cipher cipher.Block `json:"-"`
}

type LocalNode struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	WS    string `json:"ws,omitempty"` // Websocket ip
	SL    string `json:"sl,omitempty"` // Socketless pipe

	// Adapters
	websocketCache    *ristretto.Cache                        `json:"-"`
	nodeWSConnections sync.Map                                `json:"-"`
	nodes             sync.Map                                `json:"-"`
	Processors        map[string]func(*Message, string) Event `json:"-"`

	// Encryption
	Cipher cipher.Block `json:"-"`
}

var Log = log.New(log.Writer(), "pipes ", log.Flags())

func SetupCurrent(id string, token string) *LocalNode {

	if len(token) < 32 {
		panic("Token is too short (must be longer than 32 characters for AES-256 encryption)")
	}

	// Create encryption cipher
	tokenHash := sha256.Sum256([]byte(token))
	encryptionKey := tokenHash[:]

	Log.Println("Encryption key:", base64.StdEncoding.EncodeToString(encryptionKey))

	cipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		panic(err)
	}

	node := &LocalNode{
		ID:     id,
		Token:  token,
		WS:     "",
		SL:     "",
		Cipher: cipher,
	}
	node.setupCaching()
	node.setupWSStore()
	node.nodes = sync.Map{}
	node.Processors = make(map[string]func(*Message, string) Event)

	return node
}

func (n *LocalNode) SetupWS(ws string) {
	n.WS = ws
}

func (n *LocalNode) SetupSocketless(sl string) {
	n.SL = sl
}

func (local *LocalNode) GetNode(id string) *Node {

	// Get node
	obj, ok := local.nodes.Load(id)
	if !ok {
		return nil
	}
	node := obj.(Node)

	return &node
}

func (local *LocalNode) AddNode(node Node) {

	// Create encryption cipher
	tokenHash := sha256.Sum256([]byte(node.Token))
	encryptionKey := tokenHash[:]
	cipher, err := aes.NewCipher(encryptionKey)
	if err != nil {
		Log.Println("[node] Error adding node", node.ID, ":", err)
		return
	}

	node.Cipher = cipher
	local.nodes.Store(node.ID, node)
}

func (local *LocalNode) DeleteNode(node string) {
	local.nodes.Delete(node)
}

// IterateConnections iterates over all connections. If the callback returns false, the iteration stops.
func (local *LocalNode) IterateNodes(callback func(string, Node) bool) {
	local.nodes.Range(func(key, value any) bool {
		return callback(key.(string), value.(Node))
	})
}

func (local *LocalNode) ToNode() Node {
	return Node{
		ID:     local.ID,
		Token:  local.Token,
		WS:     local.WS,
		SL:     local.SL,
		Cipher: local.Cipher,
	}
}
