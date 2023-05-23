package accounts

import (
	"berith-chain/common"
	"testing"
)

func TestHex(t *testing.T) {
	var hash = common.Hash{189, 216, 255, 108, 177, 160, 220, 112, 199, 182, 132, 247, 41, 21, 41, 94, 152, 195, 20, 168, 137, 85, 251, 94, 126, 87, 126, 222, 166, 93, 11, 81}
	var chainBridgeHash = common.Hash{166, 222, 34, 63, 20, 176, 92, 173, 248, 38, 82, 135, 110, 67, 199, 71, 57, 55, 234, 78, 0, 120, 122, 116, 220, 131, 0, 236, 138, 78, 247, 122}
	t.Log(hash.Hex())
	t.Log("From ChainBridge", chainBridgeHash.Hex())
}
