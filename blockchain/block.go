package blockchain

import (
	"encoding/json"
	"time"
)

type Block struct {
	Index        int      `json:"index"`
	Timestamp    int64    `json:"timestamp"`
	PrevHash     string   `json:"prev_hash"`
	Hash         string   `json:"hash"`
	Transactions []string `json:"transactions"`
}

func NewBlock(index int, prevHash string, transactions []string) *Block {
	return &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		PrevHash:     prevHash,
		Transactions: transactions,
	}
}

func (b *Block) Serialize() []byte {
	blockBytes, _ := json.Marshal(b)
	return blockBytes
}

func DeserializeBlock(blockBytes []byte) (*Block, error) {
	var block Block
	err := json.Unmarshal(blockBytes, &block)
	return &block, err
}
