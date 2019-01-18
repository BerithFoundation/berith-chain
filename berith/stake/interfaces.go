package stake

import (
	"bitbucket.org/ibizsoftware/berith-chain/common/stakesort"
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
	Copy() StakingList
	GetRRList () *stakesort.Stakelist
	EncodeRLP(w io.Writer) error
}

type StakingInfo interface {
	Address() common.Address
	Value() *big.Int
}

// message represents a message sent to a contract.
type Transaction interface {
	From() common.Address
	//FromFrontier() (common.Address, error)
	Value() *big.Int
	Staking() bool
}

type DataBase interface {
	GetValue(key string) ([]byte, error)
	PushValue(key string, value []byte) error
	Close()
}
