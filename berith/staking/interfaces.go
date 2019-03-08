package staking

import (
	"math/big"

	"bitbucket.org/ibizsoftware/berith-chain/common"
)

//StakingList list of staked accounts
type StakingList interface {
	GetInfoWithIndex(idx int) (StakingInfo, error)
	GetInfo(address common.Address) (StakingInfo, error)
	SetInfo(info StakingInfo) error
	Delete(address common.Address) error
	Encode() ([]byte, error)
	Decode(rlpData []byte) (StakingList, error)
	Copy() StakingList
	Len() int
	Finalize()
	Print()
}

type StakingInfo interface {
	Address() common.Address
	Value() *big.Int
	BlockNumber() *big.Int
}

type DataBase interface {
	GetStakingList(key string) (StakingList, error)
	Commit(key string, stakingList StakingList) error
	NewStakingList() StakingList
	Close()
}
