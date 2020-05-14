package brtapi

import "math/big"

// [BERITH]
// Structure for returning account information
type AccountInfo struct {
	Balance      *big.Int //main balance
	StakeBalance *big.Int //staking balance
}