package test

import (
	"encoding/json"
	"math/big"
	"testing"

	"bitbucket.org/ibizsoftware/berith-chain/berith/stake"
	"bitbucket.org/ibizsoftware/berith-chain/berith/stakingdb"
	"bitbucket.org/ibizsoftware/berith-chain/common"
)

func Test5(t *testing.T) {
	test := &typ5{
		storage: make(map[string]*big.Int, 0),
	}

	t.Log(test)
	t.Log(test.storage)

	key := "0x123456789"
	value := new(big.Int).Add(big.NewInt(15), big.NewInt(5))

	test.storage[key] = value

	t.Log(test)
	t.Log(test.storage)

	result, _ := json.Marshal(test.storage)
	t.Log(result)
	t.Log(string(result))

	marshaled := new(map[string]*big.Int)

	json.Unmarshal(result, marshaled)

	t.Log(marshaled)
	// t.Log(test.storage)

}

func Test6(t *testing.T) {
	db := new(stakingdb.StakingDB)
	db.CreateDB("testStakingDB")

	blockHash0 := common.BigToHash(big.NewInt(0))

	insertList, _ := stake.NewStakingMap(db, big.NewInt(0), blockHash0)

	address1 := common.HexToAddress("0x9f3022aff3d8722043c962ee93b9a90fb580d383")
	address2 := common.HexToAddress("0xdf9bb90862483563e4fca6bb4d198587815c336b")
	address3 := common.HexToAddress("0x7cdcc6875fe66e078c815647a0bc94f25f81b500")

	t.Log("MINER1 => ", address1.String(), "VALUE => 10")
	t.Log("MINER2 => ", address2.String(), "VALUE => 20")
	t.Log("MINER3 => ", address3.String(), "VALUE => 30")

	insertList.Set(address1, big.NewInt(10))
	insertList.Set(address2, big.NewInt(20))
	insertList.Set(address3, big.NewInt(30))

	blockHash1 := common.BigToHash(big.NewInt(1))
	insertList.Commit(db, big.NewInt(1), blockHash1)

	var stakingList stake.StakingList

	stakingList, _ = stake.NewStakingMap(db, big.NewInt(1), blockHash1)

	miner1, _ := stakingList.GetMiner(0)
	miner2, _ := stakingList.GetMiner(1)
	miner3, _ := stakingList.GetMiner(2)

	t.Log("RANK1 => ", miner1.String())
	t.Log("RANK2 => ", miner2.String())
	t.Log("RANK3 => ", miner3.String())

	prev, _ := stakingList.PrevMiner(miner2)
	next, _ := stakingList.NextMiner(miner2)

	t.Log("PREV(2-1) ===>>>> ", prev.String())
	t.Log("NEXT(2+1) ===>>>> ", next.String())
}

type typ5 struct {
	storage map[string]*big.Int
}
