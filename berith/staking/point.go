/**
[BERITH]
- To calculate Selection Point
- Formula when block creation time is 10 seconds
- When generating a block every 10 seconds, 3600000 blocks are generated per year.
**/

package staking

import (
	"berith-chain/common"
	"math/big"
)

const (
	blockYear = 3600000 // When generating a block every 10 seconds, 3600000 blocks are generated per year.
)

func CalcPointBigint(prevStake, addStake, nowBlock, stakeBlock *big.Int, period uint64) *big.Int {
	correctionValue := float64(period) / common.DefaultBlockCreationSec // Value for correction when block creation time is different from the standard
	referenceBlock := int64(blockYear / correctionValue)

	ratio := new(big.Int).Mul(nowBlock, big.NewInt(100))
	ratio.Div(ratio, new(big.Int).Add(big.NewInt(referenceBlock), stakeBlock))

	if ratio.Cmp(big.NewInt(100)) == 1 {
		ratio = big.NewInt(100)
	}

	//advantage := prevStake * (prevStake / (prevStake + addStake)) * ratio / 100
	temp1 := new(big.Int).Div(prevStake, new(big.Int).Add(prevStake, addStake))
	temp2 := new(big.Int).Mul(prevStake, temp1)
	temp3 := new(big.Int).Mul(temp2, ratio)
	advantage := new(big.Int).Div(temp3, big.NewInt(100))

	//selectionPoint := prevStake + advantage + addStake
	temp1 = new(big.Int).Add(prevStake, advantage)
	selectionPoint := new(big.Int).Add(temp1, addStake)

	return selectionPoint
}
