/**
[BERITH]
- To calculate Selection Point
- Formula when block creation time is 10 seconds
- When generating a block every 10 seconds, 3600000 blocks are generated per year.
**/

package staking

import (
	"math/big"
)

const (
	BLOCK_YEAR = 3600000 // When generating a block every 10 seconds, 3600000 blocks are generated per year.
)

/*
now_block : block number
pStake : Previous staking quantity
addStake : Additional staking quantity
stake_block : Previous staking block number
epoch : Creation time per block
*/
func CalcPointUint(pStake, addStake, now_block, stake_block *big.Int, period uint64) uint64 {

	b := float64(now_block.Uint64())   // block number
	p := float64(pStake.Uint64())      // Previous staking quantity
	n := float64(addStake.Uint64())    // Additional staking quantity
	s := float64(stake_block.Uint64()) // Previous staking block number

	d := float64(period) / 10 // Value for correction when block creation time is different from the standard

	bb := BLOCK_YEAR / d // Reference block

	ratio := (b * 100) / (bb + s) // 100 is decimal point processing

	if ratio > 100 {
		ratio = 100
	}
	adv := p * ((p / (p + n)) * ratio) / 100
	result := p + adv + n

	return uint64(result)
}

func CalcPointBigint(pStake, addStake, now_block, stake_block *big.Int, period uint64) *big.Int {
	b := now_block   // block number
	p := pStake      // Previous staking quantity
	n := addStake    // Additional staking quantity
	s := stake_block // Previous staking block number

	d := float64(period) / 10 // Value for correction when block creation time is different from the standard

	bb := int64(BLOCK_YEAR / d) // Reference block

	ratio := new(big.Int).Mul(b, big.NewInt(100))
	ratio.Div(ratio, new(big.Int).Add(big.NewInt(bb), s))

	/*
		if ratio > 100 {
			ratio = 100
		}
	*/
	if ratio.Cmp(big.NewInt(100)) == 1 {
		ratio = big.NewInt(100)
	}

	//adv := p * ((p / (p + n)) * ratio) / 100
	temp1 := new(big.Int).Div(p, new(big.Int).Add(p, n))
	temp2 := new(big.Int).Mul(p, temp1)
	temp3 := new(big.Int).Mul(temp2, ratio)
	adv := new(big.Int).Div(temp3, big.NewInt(100))

	//result := p + adv + n
	r1 := new(big.Int).Add(p, adv)
	r2 := new(big.Int).Add(r1, n)

	return r2
}
