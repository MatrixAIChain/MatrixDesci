package transaction

import (
	"crypto/sha256"
	"fmt"
	"matrix-blockchain/utils"
)

// Block represents a single block in the blockchain.
type Block struct {
	Index        int           // Position of the block in the blockchain
	PreviousHash string        // Hash of the previous block
	Timestamp    int64         // Timestamp of block creation
	Transactions []Transaction // List of transactions in the block
	Hash         string        // Hash of this block
	Validator    string        // Validator who created this block
	Signature    string        // Validator's signature for the block
}

// ValidateBlock ensures that the block adheres to blockchain rules.
func ValidateBlock(newBlock, previousBlock Block) error {
	// Check if the block index is valid
	if newBlock.Index != previousBlock.Index+1 {
		return fmt.Errorf("invalid index: got %d, expected %d", newBlock.Index, previousBlock.Index+1)
	}

	// Check if the previous hash matches
	if newBlock.PreviousHash != previousBlock.Hash {
		return fmt.Errorf("invalid previous hash: got %s, expected %s", newBlock.PreviousHash, previousBlock.Hash)
	}

	// Recalculate the hash of the new block and compare
	calculatedHash := calculateBlockHash(newBlock)
	if newBlock.Hash != calculatedHash {
		return fmt.Errorf("invalid block hash: got %s, expected %s", newBlock.Hash, calculatedHash)
	}

	// Validate all transactions in the block
	for _, tx := range newBlock.Transactions {
		if !tx.Verify() {
			return fmt.Errorf("invalid transaction in block: %v", tx)
		}
	}

	// (Optional) Validate the validator's signature
	if !utils.ValidateAddress(newBlock.Validator) {
		return fmt.Errorf("invalid validator address: %s", newBlock.Validator)
	}

	return nil
}

// calculateBlockHash calculates the hash for a block.
func calculateBlockHash(block Block) string {
	data := fmt.Sprintf("%d%s%d%s%s", block.Index, block.PreviousHash, block.Timestamp, block.Transactions, block.Validator)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// IsValidChain verifies the integrity of the entire blockchain.
func IsValidChain(blockchain []Block) error {
	for i := 1; i < len(blockchain); i++ {
		err := ValidateBlock(blockchain[i], blockchain[i-1])
		if err != nil {
			return fmt.Errorf("blockchain validation failed at block %d: %v", i, err)
		}
	}
	return nil
}
