package common

import "math/big"

const (
	DefaultBlockCreationSec = 10     // Blocks are created every 10 seconds by default.
	CleanCycle              = 100000 // Standard number of levelDB clean cycle
)

var UnitForBer = big.NewInt(1e+18) // Unit to make wei to ber
