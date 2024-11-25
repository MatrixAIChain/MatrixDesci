package staking

import (
	"errors"
	"fmt"
	"sort"
)

// Validator represents a validator in the DPoS system.
type Validator struct {
	ID           string           // Validator's unique identifier (address)
	StakedAmount int64            // Total tokens staked by this validator
	Delegators   map[string]int64 // Map of delegators and their staked amounts
}

// StakingSystem manages staking operations.
type StakingSystem struct {
	Validators    map[string]*Validator // Active validators
	MaxValidators int                   // Maximum number of validators
}

// NewStakingSystem initializes a staking system.
func NewStakingSystem(maxValidators int) *StakingSystem {
	return &StakingSystem{
		Validators:    make(map[string]*Validator),
		MaxValidators: maxValidators,
	}
}

// Stake allows a user to delegate tokens to a validator.
func (s *StakingSystem) Stake(delegator string, validatorID string, amount int64) error {
	if amount <= 0 {
		return errors.New("stake amount must be positive")
	}

	validator, exists := s.Validators[validatorID]
	if !exists {
		if len(s.Validators) >= s.MaxValidators {
			return errors.New("maximum number of validators reached")
		}

		// Register a new validator
		validator = &Validator{
			ID:           validatorID,
			StakedAmount: 0,
			Delegators:   make(map[string]int64),
		}
		s.Validators[validatorID] = validator
	}

	// Add stake
	validator.StakedAmount += amount
	validator.Delegators[delegator] += amount
	return nil
}

// Unstake allows a user to withdraw their stake from a validator.
func (s *StakingSystem) Unstake(delegator string, validatorID string, amount int64) error {
	validator, exists := s.Validators[validatorID]
	if !exists {
		return errors.New("validator does not exist")
	}

	stakedAmount, exists := validator.Delegators[delegator]
	if !exists || stakedAmount < amount {
		return errors.New("not enough stake to withdraw")
	}

	// Reduce stake
	validator.StakedAmount -= amount
	validator.Delegators[delegator] -= amount

	// Remove delegator if their stake becomes zero
	if validator.Delegators[delegator] == 0 {
		delete(validator.Delegators, delegator)
	}

	// Remove validator if their total stake becomes zero
	if validator.StakedAmount == 0 {
		delete(s.Validators, validatorID)
	}

	return nil
}

// GetTopValidators returns the top N validators by stake.
func (s *StakingSystem) GetTopValidators() []*Validator {
	var validators []*Validator
	for _, v := range s.Validators {
		validators = append(validators, v)
	}

	// Sort validators by staked amount in descending order
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].StakedAmount > validators[j].StakedAmount
	})

	return validators[:s.MaxValidators]
}

// DistributeRewards distributes rewards to validators and their delegators.
func (s *StakingSystem) DistributeRewards(totalReward int64) {
	for _, validator := range s.Validators {
		// Calculate validator's share based on their stake
		reward := totalReward * validator.StakedAmount / s.TotalStake()

		// Split the reward: 25% burn, 25% research fund, 50% validator and delegators
		burn := reward / 4
		researchFund := reward / 4
		validatorShare := reward / 2

		// Validator takes their share
		validatorReward := validatorShare / 2
		validator.StakedAmount += validatorReward

		// Distribute the remaining validator share to delegators proportionally
		delegatorReward := validatorShare / 2
		for delegator, stake := range validator.Delegators {
			share := delegatorReward * stake / validator.StakedAmount
			validator.Delegators[delegator] += share
		}

		fmt.Printf("Validator %s received %d reward (Burned: %d, Research Fund: %d)\n", validator.ID, validatorReward, burn, researchFund)
	}
}

// TotalStake calculates the total stake in the system.
func (s *StakingSystem) TotalStake() int64 {
	var total int64
	for _, validator := range s.Validators {
		total += validator.StakedAmount
	}
	return total
}
