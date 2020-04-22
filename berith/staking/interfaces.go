package staking

import (
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/core/types"
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
스테이킹 정보를 관리 하는 인터페이스
*/
type StakingInfo interface {
	Address() common.Address
	Value() *big.Int
	BlockNumber() *big.Int
}

/*
[BERITH]
스테이킹 리스트 데이터 베이스 인터페이스
*/
type DataBase interface {
	GetStakers(key string) (Stakers, error)
	Commit(key string, stks Stakers) error
	NewStakers() Stakers
	Close()
	Clean(chain consensus.ChainReader, header *types.Header) error
}
