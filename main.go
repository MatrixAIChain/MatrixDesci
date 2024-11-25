package main

import (
	"fmt"
	"log"
	"matrix-blockchain/blockchain"
	"matrix-blockchain/network"
	"matrix-blockchain/staking"
	"matrix-blockchain/transaction"
)

func main() {
	// Initialize the Blockchain Database
	db, err := blockchain.OpenDatabase()
	if err != nil {
		log.Fatalf("Failed to open blockchain database: %v", err)
	}
	defer db.Close()

	// Load the latest block or initialize the genesis block if none exists
	latestBlock, err := db.GetLatestBlock()
	if err != nil {
		fmt.Println("No blocks found. Creating the genesis block...")

		// Create and save the genesis block
		genesisBlock := blockchain.createGenesisBlock()
		err = db.SaveBlock(genesisBlock)
		if err != nil {
			log.Fatalf("Failed to save genesis block: %v", err)
		}
		fmt.Println("Genesis block created successfully!")
		latestBlock = genesisBlock
	} else {
		fmt.Printf("Latest block found: %s\n", latestBlock.Hash)
	}

	// Initialize Validators
	validators := &staking.Validators{}
	validators.AddValidator("Validator1", 1000)
	validators.AddValidator("Validator2", 2000)

	// Example: Create and verify a transaction
	tx, err := transaction.NewTransaction("MRX-SenderAddress", "MRX-ReceiverAddress", 500, "PrivateKeyExample")
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	if tx.Verify() {
		fmt.Println("Transaction verified successfully.")
	} else {
		fmt.Println("Transaction verification failed.")
		return
	}

	// Example: Add a new block
	newBlock := blockchain.CreateNewBlock(latestBlock)
	newBlock.AddTransaction(tx)
	err = db.SaveBlock(newBlock)
	if err != nil {
		log.Fatalf("Failed to save the new block: %v", err)
	}
	fmt.Printf("New block added: %s\n", newBlock.Hash)

	// Example: Apply tax on a transaction
	transaction.TransactionWithTax("MRX-Address1", "MRX-Address2", 1000, 5)

	// Initialize Consensus and Voting
	consensus := blockchain.NewConsensus(validators.GetTopValidators(), newBlock.Hash)
	consensus.StartVoting()
	if consensus.IsConsensusAchieved() {
		consensus.FinalizeBlock()
		fmt.Println("Consensus achieved and block finalized.")
	} else {
		fmt.Println("Consensus not achieved. Block not finalized.")
	}

	// Distribute Rewards
	reward := staking.DistributeValidatorReward(1000)
	fmt.Println("Validator reward distribution:", reward)

	// Start the P2P Network
	p2pNetwork := network.NewP2PNetwork()
	err = p2pNetwork.Start("8080")
	if err != nil {
		log.Fatalf("Failed to start P2P network: %v", err)
	}

	// Example: Broadcast the block
	p2pNetwork.BroadcastBlock(newBlock)
	fmt.Println("Block broadcasted via P2P network.")
}
