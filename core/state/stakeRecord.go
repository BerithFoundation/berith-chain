package state

import (
	"math/big"
	"sort"
)

type StakeRecord map[*big.Int]*big.Int

func (s *StakeRecord) GetSortedKey() []*big.Int {
	var keys []*big.Int
	for k, _ := range *s {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Cmp(keys[j]) < 0 })
	return keys
}

func (s *StakeRecord) AddBalance(amount, blockNumber *big.Int) {
	(*s)[blockNumber] = new(big.Int).Add((*s)[blockNumber], amount)
}

func (s *StakeRecord) SubBalance(amount, blockNumber *big.Int) *big.Int {
	var change *big.Int
	if (*s)[blockNumber].Cmp(amount) <= 0 {
		(*s)[blockNumber] = big.NewInt(0)
		change = new(big.Int).Sub(amount, (*s)[blockNumber])
	} else {
		(*s)[blockNumber] = new(big.Int).Sub((*s)[blockNumber], amount)
	}
	return change
}
