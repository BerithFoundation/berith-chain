package staking

import (
	"fmt"
	"math/big"
	"testing"
)

func TestCalcPoint(t *testing.T) {
	add_stake := big.NewInt(1000000)
	prev_stake := big.NewInt(10000000)
	new_block := big.NewInt(7200021)
	stake_block := big.NewInt(20)
	perioid := 10
	result := CalcPoint(prev_stake, add_stake, new_block, stake_block, perioid)

	fmt.Println(result)
}
