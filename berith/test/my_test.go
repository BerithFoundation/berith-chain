package test

import (
	"fmt"
	"math/big"
	"testing"

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
	db := new(stakingdb.StakingDB)
	db.CreateDB("/Users/swk/clique/berith/stakingDB", staking.Decode, staking.Encode, staking.New)

	iter := db.Iterator()

	if iter.First() {

		for {
			key := string(iter.Key())

			fmt.Println("KEY : ", key)

			list, err := db.GetStakingList(key)

			if err != nil {
				fmt.Println("ERROR : ", err)
			} else {
				list.Print()
			}

			if !iter.Next() {
				break
			}

		}
	}

}
