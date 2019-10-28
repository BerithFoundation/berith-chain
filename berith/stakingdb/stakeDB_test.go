package stakingdb

import (
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

	stks = db.NewStakers()

	printStakers(stks.AsList())

	stks, err := db.GetStakers("test")

	if err != nil {
		println(err.Error())
	}

	printStakers(stks.AsList())

}

func printStakers(list []common.Address) {
	println("====================[STAKERS]======================")
	for i, v := range list {
		println("[", i, ",", v.Hex(), "]")
	}
}
