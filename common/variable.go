package common

import "math/big"

const (
	DefaultBlockCreationSec = 10        // Blocks are created every 10 seconds by default.
	CleanCycle              = 100000    // Standard number of levelDB clean cycle

	KeyForGCMode            = "gcmode"  // Key to get user selected gc mode from walletDB
	KeyForGCModeChangYn     = "changeYn"
	GCModeArchive           = "archive"
	GCModeFull              = "full"
)

var UnitForBer = big.NewInt(1e+18) // Unit to make wei to ber
