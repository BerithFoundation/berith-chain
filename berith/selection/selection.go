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
	for _, stk := range list {
		point := state.GetPoint(stk).Uint64()
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