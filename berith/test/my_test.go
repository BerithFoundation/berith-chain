package test

import (
	"testing"

	"github.com/BerithFoundation/berith-chain/rlp"

	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/core/state"
	"github.com/BerithFoundation/berith-chain/trie"
)

func Test01(t *testing.T) {
	db, err := berithdb.NewLDBDatabase("C:\\Users\\ibizsoftware\\workspace\\gopath\\bin\\pos04\\berith\\chaindata", 0, 0)

	if err != nil {
		println("ERR : ", err.Error())
		return
	}

	root := common.HexToHash("0x9918f62ffef79179d9cddc97d18b8e06f0941564fb73dcbe6b1f288c35d7377a")

	stateDB := state.NewDatabaseWithCache(db, 2048)

	tri, err := trie.NewSecure(root, stateDB.TrieDB(), 0)

	if err != nil {
		println("ERR : ", err.Error())
		return
	}

	addr := common.HexToAddress("0x78c2b0dfde452677ccd0cd00465e7cca0e3c5353")

	enc, err := tri.TryGet(addr[:])

	if err != nil {
		println("ERR : ", err.Error())
		return
	}

	account := state.Account{}

	rlp.DecodeBytes(enc, &account)

	println("MAIN : ", account.Balance.String())

}

// import (
// 	"fmt"
// 	"math/big"
// 	"testing"

// 	"github.com/BerithFoundation/berith-chain/rlp"

// 	lru "github.com/hashicorp/golang-lru"

// 	"github.com/BerithFoundation/berith-chain/berith/staking"
// 	"github.com/BerithFoundation/berith-chain/berith/stakingdb"
// 	"github.com/BerithFoundation/berith-chain/common"
// )

// type StakingInfo struct {
// 	address     common.Address
// 	value       *big.Int
// 	blockNumber *big.Int
// 	reward      *big.Int
// }

// func (stk StakingInfo) Address() common.Address { return stk.address }
// func (stk StakingInfo) Value() *big.Int         { return stk.value }
// func (stk StakingInfo) BlockNumber() *big.Int   { return stk.blockNumber }
// func (stk StakingInfo) Reward() *big.Int        { return stk.reward }

// func Test1(t *testing.T) {
// 	cache, _ := lru.NewARC(10)
// 	db := new(stakingdb.StakingDB)
// 	db.CreateDB("C:\\Users\\ibizsoftware\\testdb\\db", staking.Decode, staking.Encode, staking.New)
// 	list := db.NewStakingList()
// 	for i := int64(0); i < 5; i++ {
// 		val, _ := new(big.Int).SetString("100000000000000000000000", 10)
// 		list.SetInfo(StakingInfo{
// 			address:     common.BigToAddress(big.NewInt(i)),
// 			value:       val,
// 			blockNumber: big.NewInt(10),
// 			reward:      big.NewInt(0),
// 		})
// 	}
// 	list.Sort()
// 	diff, rank, reordered := list.GetDifficultyAndRank(common.BigToAddress(big.NewInt(1)), 10, 10)
// 	fmt.Println(diff.String(), rank, reordered)
// 	list.Print()
// 	encoded, _ := rlp.EncodeToBytes(list)
// 	cache.Add("stk", encoded)
// 	db.Commit("stk", list)

// 	incache, _ := cache.Get("stk")

// 	outcache, ok := incache.([]byte)

// 	if !ok {
// 		fmt.Println("failed to get list in cache")
// 	}

// 	list1, err := staking.Decode(outcache)

// 	if err != nil {
// 		fmt.Println("failed to decode data in cache", err)
// 	}

// 	list1.Print()

// 	list2, _ := db.GetStakingList("stk")

// 	list2.Print()

// }
