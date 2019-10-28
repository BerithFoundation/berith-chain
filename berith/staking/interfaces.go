package staking

import (
	"io"
	"math/big"

	"github.com/BerithFoundation/berith-chain/rlp"

	"github.com/BerithFoundation/berith-chain/common"
)

//StakingList list of staked accounts
/*
[BERITH]
스테이킹 리스트 인터페이스
*/
// type StakingList interface {
// 	GetInfo(address common.Address) (StakingInfo, error)
// 	SetInfo(info StakingInfo) error
// 	Delete(address common.Address) error
// 	Encode() ([]byte, error)
// 	Decode(rlpData []byte) (StakingList, error)
// 	Copy() StakingList
// 	Len() int
// 	Print()
// 	GetJoinRatio(address common.Address, blockNumber uint64, states *state.StateDB) float64
// 	Sort()
// 	ClearTable()
// 	GetDifficultyAndRank(addr common.Address, blockNumber uint64, states *state.StateDB) (*big.Int, int, bool)
// 	ToArray() []common.Address
// }

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
}
