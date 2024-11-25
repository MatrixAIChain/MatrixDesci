package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Block represents a block in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Validator    string
}

// Transaction represents a blockchain transaction
type Transaction struct {
	From   string
	To     string
	Amount int
}

// Blockchain represents a simple blockchain with a LevelDB backend
type Blockchain struct {
	db *leveldb.DB
}

// NewBlockchain initializes a new Blockchain with LevelDB
func NewBlockchain(dbPath string) (*Blockchain, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return &Blockchain{db: db}, nil
}

// AddBlock adds a block to the blockchain (LevelDB storage)
func (bc *Blockchain) AddBlock(block *Block) error {
	// Serialize the block and store it in the database
	blockBytes := block.Serialize()
	err := bc.db.Put([]byte(fmt.Sprintf("block_%d", block.Index)), blockBytes, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetBlock retrieves a block by its index
func (bc *Blockchain) GetBlock(index int) (*Block, error) {
	blockBytes, err := bc.db.Get([]byte(fmt.Sprintf("block_%d", index)), nil)
	if err != nil {
		return nil, err
	}

	block, err := DeserializeBlock(blockBytes)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetLatestBlock retrieves the latest block in the blockchain
func (bc *Blockchain) GetLatestBlock() (*Block, error) {
	iter := bc.db.NewIterator(util.BytesPrefix([]byte("block_")), nil)
	defer iter.Release()

	var latestBlock *Block
	for iter.Next() {
		blockBytes := iter.Value()
		block, err := DeserializeBlock(blockBytes)
		if err != nil {
			return nil, err
		}
		latestBlock = block
	}

	if latestBlock == nil {
		return nil, fmt.Errorf("no blocks found")
	}

	return latestBlock, nil
}

// Serialize converts a block into bytes for storage
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	// Serialize the block into a byte slice (this is a simplistic example)
	// You may want to use encoding such as JSON, Protocol Buffers, or other methods for better performance
	buffer.WriteString(fmt.Sprintf("%d%s%s", block.Index, block.Timestamp, block.PrevHash))
	for _, tx := range block.Transactions {
		buffer.WriteString(fmt.Sprintf("%s%s%d", tx.From, tx.To, tx.Amount))
	}
	return buffer.Bytes()
}

// DeserializeBlock converts bytes back into a Block
func DeserializeBlock(data []byte) (*Block, error) {
	// This is a simplistic deserialization (in a real case, you might use a more structured deserialization)
	// You should parse the data and reconstruct the Block here.
	block := &Block{}
	// Deserialization logic should happen here.
	return block, nil
}

// CreateGenesisBlock creates the first block (genesis block) of the blockchain
func CreateGenesisBlock() *Block {
	tx := Transaction{
		From:   "GENESIS",
		To:     "MRX-InitialWallet",
		Amount: 500000000, // Pre-mined supply of 500 million
	}

	block := &Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Transactions: []Transaction{tx},
		PrevHash:     "",
		Hash:         calculateHash(0, "", []Transaction{tx}),
		Validator:    "GENESIS_VALIDATOR",
	}

	return block
}

// calculateHash generates a SHA256 hash for a block
func calculateHash(index int, prevHash string, txs []Transaction) string {
	hashInput := fmt.Sprintf("%d%s%v", index, prevHash, txs)
	hash := sha256.Sum256([]byte(hashInput))
	return fmt.Sprintf("%x", hash)
}
