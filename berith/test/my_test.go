package test

import (
	"math/big"
	"testing"

	"bitbucket.org/ibizsoftware/berith-chain/berith/staking"
	"bitbucket.org/ibizsoftware/berith-chain/berith/stakingdb"
	"bitbucket.org/ibizsoftware/berith-chain/common"
)

type StakingInfo struct {
	address     common.Address
	value       *big.Int
	blockNumber *big.Int
}

func (stk StakingInfo) Address() common.Address { return stk.address }
func (stk StakingInfo) Value() *big.Int         { return stk.value }
func (stk StakingInfo) BlockNumber() *big.Int   { return stk.blockNumber }

func Test1(t *testing.T) {
	db := new(stakingdb.StakingDB)
	db.CreateDB("/Users/swk/clique/geth/stakingDB", staking.Decode, staking.Encode, staking.New)

	stakingMap := db.NewStakingList()

	stakingMap.SetInfo(StakingInfo{
		address:     common.Address{},
		value:       big.NewInt(4000),
		blockNumber: big.NewInt(10),
	})

	stakingMap.Print()

	encoding, err := stakingMap.Encode()

	if err != nil {
		t.Error(err)
	}

	var decoding staking.StakingList

	decoding, err = staking.Decode(encoding)

	if err != nil {
		t.Error(err)
	}

	decoding.Print()

	//t.Log(common.BytesToAddress(rlpVal).Hex())

}
