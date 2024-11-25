package staking

import (
	"fmt"
)

type Reward struct {
	TotalAmount        int64
	BurnAmount         int64
	ValidatorAmount    int64
	ResearchFundAmount int64
}

func DistributeValidatorReward(totalRewards int64) Reward {
	burn := totalRewards * 25 / 100
	researchFund := totalRewards * 25 / 100
	validator := totalRewards * 50 / 100

	reward := Reward{
		TotalAmount:        totalRewards,
		BurnAmount:         burn,
		ResearchFundAmount: researchFund,
		ValidatorAmount:    validator,
	}

	fmt.Printf("Distributed %d to validator, %d burned, %d to research fund\n", validator, burn, researchFund)

	return reward
}
