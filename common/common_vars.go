package common

import "math/big"

const (
	DefaultBlockCreationSec = 10 // Blocks are created every 10 seconds by default.
)

var UnitForBer = big.NewInt(1e+18) // Unit to make wei to ber