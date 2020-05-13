package staking

import (
	"io"
	"math/big"

	"github.com/BerithFoundation/berith-chain/rlp"

	"github.com/BerithFoundation/berith-chain/common"
)

type Stakers interface {
	Put(common.Address)
	Remove(common.Address)
	IsContain(common.Address) bool
	AsList() []common.Address
	FetchFromList([]common.Address)
	EncodeRLP(io.Writer) error
	DecodeRLP(*rlp.Stream) error
}

/*
[BERITH]
Interface to manage staking information
*/
type StakingInfo interface {
	Address() common.Address
	Value() *big.Int
	BlockNumber() *big.Int
}

/*
[BERITH]
Interface for stakingDB
*/
type DataBase interface {
	GetStakers(key string) (Stakers, error)
	Commit(key string, stks Stakers) error
	NewStakers() Stakers
	Close()
}
