package stakingdb

import (
	"fmt"
	"testing"

	"github.com/BerithFoundation/berith-chain/berith/staking"
	"github.com/BerithFoundation/berith-chain/common"
)

func Test01(t *testing.T) {

	db := new(StakingDB)
	db.CreateDB("C:/Users/ibizsoftware/test/stakingdb/", staking.NewStakers)

	stks := db.NewStakers()
	addr1 := common.BytesToAddress([]byte("1"))
	addr2 := common.BytesToAddress([]byte("2"))

	printStakers(stks.AsList())

	stks.Put(addr1)
	stks.Put(addr2)

	printStakers(stks.AsList())

	stks.Remove(addr1)

	printStakers(stks.AsList())

	stks.Remove(addr1)

	printStakers(stks.AsList())

	stks.Put(addr1)
	stks.Put(addr1)

	printStakers(stks.AsList())

	db.Commit("test", stks)

}

func printStakers(list []common.Address) {
	println("====================[STAKERS]======================")
	for i, v := range list {
		fmt.Printf("[%d,%s]\n", i, v.Hex())
	}
}
