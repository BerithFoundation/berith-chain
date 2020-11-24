package selection

import "github.com/BerithFoundation/berith-chain/common"

/*
[BERITH]
Structure that holds the information of the staking accounts for election
*/
type Candidate struct {
	address common.Address // Account address
	point   uint64         // Points in the account (probability of being drawn: my points / points of all users)
	val     uint64         // Value used to elect block constructor
}

type JSONCandidate struct {
	Address string `json:"address"`
	Point   uint64 `json:"point"`
	Value   uint64 `json:"value"`
}

func (c *Candidate) GetPoint() uint64 {
	return c.point
}
