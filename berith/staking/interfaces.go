package staking

import (
	"math/big"

	"github.com/BerithFoundation/berith-chain/common"
)

//StakingList list of staked accounts
/*
[BERITH]
스테이킹 리스트 인터페이스
*/
type StakingList interface {
	GetInfo(address common.Address) (StakingInfo, error)
	SetInfo(info StakingInfo) error
	Delete(address common.Address) error
	Encode() ([]byte, error)
	Decode(rlpData []byte) (StakingList, error)
	Copy() StakingList
	Len() int
	Print()
	GetJoinRatio(address common.Address, blockNumber, period uint64) float64
	Sort()
	ClearTable()
	GetDifficultyAndRank(addr common.Address, blockNumber, period uint64) (*big.Int, int, bool)
	ToArray() []common.Address
}

/*
[BERITH]
스테이킹 정보를 관리 하는 인터페이스
*/
type StakingInfo interface {
	Address() common.Address
	Value() *big.Int
	BlockNumber() *big.Int
	Reward() *big.Int
}

/*
[BERITH]
스테이킹 리스트 데이터 베이스 인터페이스
*/
type DataBase interface {
	GetStakingList(key string) (StakingList, error)
	Commit(key string, stakingList StakingList) error
	NewStakingList() StakingList
	Close()
}
