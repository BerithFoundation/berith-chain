package staking

import (
	"math/big"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/consensus"
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
	Vote(chain consensus.ChainReader, number uint64, hash common.Hash, epoch uint64, period uint64)
	Print()
	GetRoundJoinRatio() *map[common.Address]int
	SetMiner(address common.Address)
	InitMiner()
	GetMiners() map[common.Address]bool
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
