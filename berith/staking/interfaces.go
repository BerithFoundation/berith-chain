package staking

import (
	"math/big"

	"github.com/BerithFoundation/berith-chain/common"
)

//StakingList list of staked accounts
type StakingList interface {
	SetTarget(target common.Hash)
	GetTarget() common.Hash
	GetInfoWithIndex(idx int) (StakingInfo, error)
	GetInfo(address common.Address) (StakingInfo, error)
	SetInfo(info StakingInfo) error
	Delete(address common.Address) error
	Encode() ([]byte, error)
	Decode(rlpData []byte) (StakingList, error)
	Copy() StakingList
	Len() int
	Print()
	GetRoi(address common.Address) float64
	SetMiner(address common.Address)
	InitMiner()
	GetMiners() map[common.Address]bool
	Sort()
	ClearTable()
	GetDifficulty(addr common.Address, blockNumber, period uint64) (*big.Int, bool)
	ToArray() []common.Address
}

type StakingInfo interface {
	Address() common.Address
	Value() *big.Int
	BlockNumber() *big.Int
	Reward() *big.Int
}

type DataBase interface {
	GetStakingList(key string) (StakingList, error)
	Commit(key string, stakingList StakingList) error
	NewStakingList() StakingList
	Close()
}
