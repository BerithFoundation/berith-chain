package test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/BerithFoundation/berith-chain/rlp"

	lru "github.com/hashicorp/golang-lru"

	"github.com/BerithFoundation/berith-chain/berith/staking"
	"github.com/BerithFoundation/berith-chain/berith/stakingdb"
	"github.com/BerithFoundation/berith-chain/common"
)

type StakingInfo struct {
	address     common.Address
	value       *big.Int
	blockNumber *big.Int
	reward      *big.Int
}

func (stk StakingInfo) Address() common.Address { return stk.address }
func (stk StakingInfo) Value() *big.Int         { return stk.value }
func (stk StakingInfo) BlockNumber() *big.Int   { return stk.blockNumber }
func (stk StakingInfo) Reward() *big.Int        { return stk.reward }

func Test1(t *testing.T) {
	cache, _ := lru.NewARC(10)
	db := new(stakingdb.StakingDB)
	db.CreateDB("C:\\Users\\ibizsoftware\\testdb\\db", staking.Decode, staking.Encode, staking.New)
	list := db.NewStakingList()
	for i := int64(0); i < 5; i++ {
		val, _ := new(big.Int).SetString("100000000000000000000000", 10)
		list.SetInfo(StakingInfo{
			address:     common.BigToAddress(big.NewInt(i)),
			value:       val,
			blockNumber: big.NewInt(10),
			reward:      big.NewInt(0),
		})
	}
	list.Sort()
	diff, rank, reordered := list.GetDifficultyAndRank(common.BigToAddress(big.NewInt(1)), 10, 10)
	fmt.Println(diff.String(), rank, reordered)
	list.Print()
	encoded, _ := rlp.EncodeToBytes(list)
	cache.Add("stk", encoded)
	db.Commit("stk", list)

	incache, _ := cache.Get("stk")

	outcache, ok := incache.([]byte)

	if !ok {
		fmt.Println("failed to get list in cache")
	}

	list1, err := staking.Decode(outcache)

	if err != nil {
		fmt.Println("failed to decode data in cache", err)
	}

	list1.Print()

	list2, _ := db.GetStakingList("stk")

	list2.Print()

}
