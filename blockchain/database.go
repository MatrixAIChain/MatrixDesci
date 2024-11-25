package blockchain

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const (
	blockchainDBFile = "blockchain.db"
	blocksBucket     = "blocks"
	latestBlockKey   = "latest"
)

// Database represents the blockchain database.
type Database struct {
	db *bolt.DB
}

// OpenDatabase opens or creates the blockchain database.
func OpenDatabase() (*Database, error) {
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create the blocks bucket if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(blocksBucket))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create blocks bucket: %v", err)
	}

	return &Database{db: db}, nil
}

// SaveBlock stores a block in the database.
func (db *Database) SaveBlock(block *Block) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			return fmt.Errorf("blocks bucket not found")
		}

		data, err := json.Marshal(block)
		if err != nil {
			return fmt.Errorf("failed to serialize block: %v", err)
		}

		err = bucket.Put([]byte(block.Hash), data)
		if err != nil {
			return fmt.Errorf("failed to save block: %v", err)
		}

		// Update the latest block reference
		err = bucket.Put([]byte(latestBlockKey), []byte(block.Hash))
		if err != nil {
			return fmt.Errorf("failed to update latest block: %v", err)
		}

		return nil
	})
}

// GetBlock retrieves a block by its hash.
func (db *Database) GetBlock(hash string) (*Block, error) {
	var block *Block

	err := db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			return fmt.Errorf("blocks bucket not found")
		}

		data := bucket.Get([]byte(hash))
		if data == nil {
			return fmt.Errorf("block not found")
		}

		block = &Block{}
		err := json.Unmarshal(data, block)
		if err != nil {
			return fmt.Errorf("failed to deserialize block: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return block, nil
}

// GetLatestBlock retrieves the latest block in the blockchain.
func (db *Database) GetLatestBlock() (*Block, error) {
	var latestBlock *Block

	err := db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			return fmt.Errorf("blocks bucket not found")
		}

		latestHash := bucket.Get([]byte(latestBlockKey))
		if latestHash == nil {
			return fmt.Errorf("no latest block found")
		}

		data := bucket.Get(latestHash)
		if data == nil {
			return fmt.Errorf("latest block data not found")
		}

		latestBlock = &Block{}
		err := json.Unmarshal(data, latestBlock)
		if err != nil {
			return fmt.Errorf("failed to deserialize latest block: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return latestBlock, nil
}

// Close closes the database connection.
func (db *Database) Close() {
	err := db.db.Close()
	if err != nil {
		log.Printf("Failed to close database: %v", err)
	}
}
