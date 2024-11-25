package blockchain

import (
	"fmt"
	"math/rand"
	"matrix-blockchain/staking"
)

type Consensus struct {
	Validators []staking.Validator
	Quorum     int
	BlockHash  string
	Votes      map[string]bool // Map of validator ID and their vote (yes/no)
}

func NewConsensus(validators []staking.Validator, blockHash string) *Consensus {
	return &Consensus{
		Validators: validators,
		Quorum:     len(validators) * 2 / 3, // 2/3 majority for consensus
		BlockHash:  blockHash,
		Votes:      make(map[string]bool),
	}
}

func (c *Consensus) StartVoting() {
	// Randomly pick a subset of validators to vote on the block
	for _, validator := range c.Validators {
		// Simulate a vote decision based on validator's stake
		if rand.Float32() < float32(validator.StakedAmount)/10000 {
			c.Votes[validator.ID] = true
		} else {
			c.Votes[validator.ID] = false
		}
	}
}

func (c *Consensus) IsConsensusAchieved() bool {
	votesFor := 0
	for _, vote := range c.Votes {
		if vote {
			votesFor++
		}
	}
	return votesFor >= c.Quorum
}

func (c *Consensus) FinalizeBlock() bool {
	if c.IsConsensusAchieved() {
		fmt.Println("Consensus reached, block finalized.")
		return true
	}
	fmt.Println("Consensus not reached, block rejected.")
	return false
}
