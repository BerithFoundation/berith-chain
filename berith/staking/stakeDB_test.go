package staking

import (
	"errors"
	"os"
	"testing"

	"github.com/BerithFoundation/berith-chain/common"
)

func TestStakers(t *testing.T) {
	db := new(StakingDB)
	db.CreateDB(os.TempDir()+"/stakingdb/", NewStakers)

	err := errors.New("result is incorrect than expected")

	stks := db.NewStakers()
	addr1 := common.BytesToAddress([]byte("1"))
	addr2 := common.BytesToAddress([]byte("2"))

	if stks.IsContain(addr1) || stks.IsContain(addr2) || len(stks.AsList()) != 0 {
		t.Error(err)
	}

	stks.Put(addr1)
	stks.Put(addr2)

	if !stks.IsContain(addr1) || !stks.IsContain(addr2) || len(stks.AsList()) != 2 {
		t.Error(err)
	}

	stks.Remove(addr1)

	if stks.IsContain(addr1) || !stks.IsContain(addr2) || len(stks.AsList()) != 1 {
		t.Error(err)
	}

	stks.Remove(addr1)

	if stks.IsContain(addr1) || !stks.IsContain(addr2) || len(stks.AsList()) != 1 {
		t.Error(err)
	}

	stks.Put(addr1)
	stks.Put(addr1)

	if !stks.IsContain(addr1) || !stks.IsContain(addr2) || len(stks.AsList()) != 2 {
		t.Error(err)
	}

	db.Commit("test", stks)

	stks = db.NewStakers()

	if stks.IsContain(addr1) || stks.IsContain(addr2) || len(stks.AsList()) != 0 {
		t.Error(err)
	}

	stks, dbErr := db.GetStakers("test")

	if dbErr != nil {
		t.Error(dbErr)
	}

	if !stks.IsContain(addr1) || !stks.IsContain(addr2) || len(stks.AsList()) != 2 {
		t.Error(err)
	}

}
