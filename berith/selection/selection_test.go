package selection

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/BerithFoundation/berith-chain/berith/staking"
	"github.com/BerithFoundation/berith-chain/berith/stakingdb"

	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/core/state"
)

/*
[BERITH]
선출 로직 테스트
*/
func TestVoting2(t *testing.T) {

	st, _ := state.New(common.Hash{}, state.NewDatabase(berithdb.NewMemDatabase()))

	stakingDB := stakingdb.StakingDB{}
	stakingDB.CreateDB("C:/Users/ibizsoftware/test/stakingdb", staking.NewStakers)

	stks := stakingDB.NewStakers()

	blockNumber := big.NewInt(100)
	eth := big.NewInt(1e+18)
	value := new(big.Int).Mul(big.NewInt(100000), eth)
	for i := 0; i < 100; i++ {

		addr := common.BigToAddress(big.NewInt(int64(i)))

		st.AddStakeBalance(addr, value, blockNumber)
		stks.Put(addr)

		prev_stake := new(big.Int).Div(st.GetStakeBalance(addr), big.NewInt(1e+18))
		add_stake := new(big.Int).Div(value, big.NewInt(1e+18))
		now_block := blockNumber
		stake_block := new(big.Int).Set(st.GetStakeUpdated(addr))
		period := uint64(40)

		result := staking.CalcPointBigint(prev_stake, add_stake, now_block, stake_block, period)
		st.SetPoint(addr, result)
	}

	results := SelectBlockCreator(blockNumber.Uint64(), common.Hash{}, stks, st)

	m := make(map[int]common.Address)

	for k, v := range results {
		m[v.Rank] = k
	}

	for i := 1; i <= 100; i++ {
		result := results[m[i]]

		fmt.Printf("[%d,%s,%s]\n", result.Rank, m[i].Hex(), result.Score)
	}
	fmt.Println("LEN : ", len(results))
}
