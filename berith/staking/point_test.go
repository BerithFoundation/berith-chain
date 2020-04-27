package staking

import (
	"fmt"
	"math/big"
	"testing"
)

/*
[BERITH]
Election point calculation test
*/
func TestCalcPoint(t *testing.T) {
	add_stake := big.NewInt(1000000)
	prev_stake := big.NewInt(10000000)
	new_block := big.NewInt(7200021)
	stake_block := big.NewInt(20)
	perioid := uint64(10)
	result := CalcPointUint(prev_stake, add_stake, new_block, stake_block, perioid)

	fmt.Println(result)
}

/*
[BERITH]
Election point calculation test
*/
func TestCalcPoint2(t *testing.T) {
	add_stake := big.NewInt(1000000)
	prev_stake := big.NewInt(10000000)
	new_block := big.NewInt(7200021)
	stake_block := big.NewInt(20)
	perioid := uint64(10)
	result := CalcPointBigint(prev_stake, add_stake, new_block, stake_block, perioid)

	fmt.Println(result)
}
