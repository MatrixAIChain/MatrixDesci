package staking

import (
	"fmt"
)

type Validator struct {
	ID           string
	StakedAmount int64
	Reward       int64
}

type Validators struct {
	ValidatorList []Validator
}

func (v *Validators) AddValidator(id string, stakedAmount int64) {
	validator := Validator{
		ID:           id,
		StakedAmount: stakedAmount,
	}
	v.ValidatorList = append(v.ValidatorList, validator)
}

func (v *Validators) GetTopValidators() []Validator {
	// Sort or pick top N validators (e.g., by staked amount)
	return v.ValidatorList[:100]
}

func (v *Validators) DistributeRewards(totalRewards int64) {
	for _, validator := range v.GetTopValidators() {
		reward := totalRewards * 50 / 100
		burn := totalRewards * 25 / 100
		researchFund := totalRewards * 25 / 100
		fmt.Printf("Validator %s gets %d, %d burned, %d for research\n", validator.ID, reward, burn, researchFund)
		// Apply logic for burning and rewarding
	}
}
