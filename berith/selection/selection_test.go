package selection

import (
	"math/big"
	"testing"

	"github.com/BerithFoundation/berith-chain/berith/staking"

	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/core/state"
)

/*
[BERITH]
선출 로직 테스트
*/
func TestSelectBlockCreator(t *testing.T) {
	expectedResults := map[common.Address]VoteResult{
		common.HexToAddress("0000000000000000000000000000000000000001"): VoteResult{
			Score: big.NewInt(5000000),
			Rank:  1,
		},
		common.HexToAddress("0000000000000000000000000000000000000000"): VoteResult{
			Score: big.NewInt(4002000),
			Rank:  2,
		},
		common.HexToAddress("0000000000000000000000000000000000000004"): VoteResult{
			Score: big.NewInt(3004000),
			Rank:  3,
		},
		common.HexToAddress("0000000000000000000000000000000000000003"): VoteResult{
			Score: big.NewInt(2006000),
			Rank:  4,
		},
		common.HexToAddress("0000000000000000000000000000000000000002"): VoteResult{
			Score: big.NewInt(1008000),
			Rank:  5,
		},
	}

	st, _ := state.New(common.Hash{}, state.NewDatabase(berithdb.NewMemDatabase()))

	stks := staking.NewStakers()

	blockNumber := big.NewInt(100)
	eth := big.NewInt(1e+18)
	value := new(big.Int).Mul(big.NewInt(100000), eth)
	for i := 0; i < 5; i++ {

		addr := common.BigToAddress(big.NewInt(int64(i)))

		st.AddStakeBalance(addr, value, blockNumber)
		stks.Put(addr)

		prevStake := new(big.Int).Div(st.GetStakeBalance(addr), big.NewInt(1e+18))
		addStake := new(big.Int).Div(value, big.NewInt(1e+18))
		nowBlock := blockNumber
		stakeBlock := new(big.Int).Set(st.GetStakeUpdated(addr))
		period := uint64(40)

		point := staking.CalcPointBigint(prevStake, addStake, nowBlock, stakeBlock, period)
		st.SetPoint(addr, point)
	}

	results := SelectBlockCreator(blockNumber.Uint64(), common.Hash{}, stks, st)

	for addr, result := range results {
		expected, ok := expectedResults[addr]
		if !ok {
			t.Errorf("%s isn't in expected result", addr)
		}
		if expected.Rank != result.Rank || expected.Score.Cmp(result.Score) != 0 {
			t.Errorf("expected result is [%d, %s] but, [%d, %s]", expected.Rank, expected.Score.String(), result.Rank, result.Score.String())
		}
	}

	if len(results) < 5 {
		t.Errorf("only %d user selected [expected : 100]", len(results))
	}
}
