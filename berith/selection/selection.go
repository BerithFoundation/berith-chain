/**
[BERITH]
package in charge of electing operation
*/

package selection

import (
	"math/big"
	"sort"

	"github.com/BerithFoundation/berith-chain/params"

	"github.com/BerithFoundation/berith-chain/berith/staking"
	"github.com/BerithFoundation/berith-chain/core/state"

	"github.com/BerithFoundation/berith-chain/common"
)

/*
[BERITH]
Entry function to elect Block Creator
Returns the elected Block Creator map.
*/
func SelectBlockCreator(config *params.ChainConfig, number uint64, hash common.Hash, stks staking.Stakers, state *state.StateDB) VoteResults {
	result := make(VoteResults)

	// Get and Sort staker list
	list := sortableList(stks.AsList())
	if len(list) == 0 {
		return result
	}
	sort.Sort(list)

	// Make Candidates data structure
	cddts := NewCandidates()
	blockNumber := big.NewInt(int64(number))

	/*
		[Berith]
		In accordance with the addition of the Stake Balance limit, targets with a Stake Balance limit or higher are recalculated.
	*/
	for _, stk := range list {
		stakeBalance := state.GetStakeBalance(stk)
		var point uint64

		if config.IsBIP4(blockNumber) && stakeBalance.Cmp(config.Bsrr.LimitStakeBalance) == 1 {
			limitStakeBalanceInBer := new(big.Int).Div(config.Bsrr.LimitStakeBalance, common.UnitForBer)
			lastStkBlock := new(big.Int).Set(state.GetStakeUpdated(stk))
			advantage := calcAdvForExceededPoint(blockNumber, lastStkBlock, config.Bsrr.Period, common.BigIntToBigFloat(limitStakeBalanceInBer))

			point = new(big.Int).Add(limitStakeBalanceInBer, advantage).Uint64()
		} else {
			point = state.GetPoint(stk).Uint64()
		}

		cddts.Add(Candidate{
			point:   point,
			address: stk,
		})
	}

	// Call block creator function
	if config.IsBIP3(big.NewInt(int64(number))) {
		result = cddts.selectBIP3BlockCreator(config, number)
	} else {
		result = cddts.selectBlockCreator(config, number)
	}

	return result
}

/*
	[Berith]
	A function that newly calculates the elected point advantage for holders who have exceeded the Stake Balance limit
*/
func calcAdvForExceededPoint(nowBlockNumber, stakeBlockNumber *big.Int, period uint64, limitStakeBalanceInBer *big.Float) *big.Int {
	d := float64(period) / 10 //공식이 10초 단위 이기때문에 맞추기 위함 (perioid 를 제네시스로 변경하면 자동으로 변경되기 위함)

	bb := float64(staking.BlockYear / d) //기준 블록

	//ratio := (b * 100)  / (bb + s) //100은 소수점 처리
	ratio := new(big.Float).Mul(new(big.Float).SetInt(nowBlockNumber), big.NewFloat(100))
	ratio.Quo(ratio, new(big.Float).Add(big.NewFloat(bb), new(big.Float).SetInt(stakeBlockNumber)))

	/*
		if ratio > 100 {
			ratio = 100
		}
	*/
	if ratio.Cmp(big.NewFloat(100)) == 1 {
		ratio = big.NewFloat(100)
	}

	temp1 := new(big.Float).Quo(limitStakeBalanceInBer, new(big.Float).Add(limitStakeBalanceInBer, big.NewFloat(0)))
	temp2 := new(big.Float).Mul(limitStakeBalanceInBer, temp1)
	temp3 := new(big.Float).Mul(temp2, ratio)
	adv := new(big.Int)
	new(big.Float).Quo(temp3, big.NewFloat(100)).Int(adv)

	return adv
}

func GetCandidates(number uint64, hash common.Hash, stks staking.Stakers, state *state.StateDB) *JSONCandidates {
	list := sortableList(stks.AsList())

	sort.Sort(list)

	cddts := NewCandidates()
	for _, stk := range list {
		point := state.GetPoint(stk).Uint64()
		cddts.Add(Candidate{
			point:   point,
			address: stk,
		})
	}

	jsonCddt := make([]JSONCandidate, 0)
	for _, cddt := range cddts.selections {
		jsonCddt = append(jsonCddt, JSONCandidate{
			Address: cddt.address.Hex(),
			Point:   cddt.point,
			Value:   cddt.val,
		})
	}
	return &JSONCandidates{
		User:  jsonCddt,
		Total: cddts.total,
	}
}
