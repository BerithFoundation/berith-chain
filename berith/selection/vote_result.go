package selection

import (
	"berith-chain/common"
	"math/big"
)

/*
[Berith]
Object that stores voting results for election
*/
type VoteResult struct {
	Score *big.Int `json:"score"`
	Rank  int      `json:"rank"`
}

type VoteResults map[common.Address]VoteResult