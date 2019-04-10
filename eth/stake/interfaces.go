package stake

import (
	"io"
	"math/big"

	"bitbucket.org/ibizsoftware/berith-chain/common"
)

//stakingListKey trieDB's key for staking list
const stakingListKey = "staking_list"

//StakingList list of staked accounts
type StakingList interface {
	Get(address common.Address) (StakingInfo, error)
	Set(address common.Address, x interface{}) error
	Delete(address common.Address) error
	EncodeRLP(w io.Writer) error
	Commit(db DataBase, blockNumber *big.Int, blockHash common.Hash) error
	NextMiner(address common.Address) (common.Address, error)
	PrevMiner(address common.Address) (common.Address, error)
	GetMiner(index int) (common.Address, error)
	Len() int
}

type StakingInfo interface {
	Address() common.Address
	Value() *big.Int
}

type DataBase interface {
	GetValue(key string) ([]byte, error)
	PushValue(key string, value []byte) error
	Close()
}
