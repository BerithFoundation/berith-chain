package common

import (
	"fmt"
	"testing"
)

func TestAddressHex(t *testing.T) {
	// var coinbase = Address{88, 179, 149, 164, 186, 25, 94, 3, 68, 90, 6, 92, 182, 180, 241, 149, 87, 45, 251, 231}
	// var msgSender = Address{163, 162, 106, 43, 205, 33, 11, 249, 102, 247, 149, 58, 188, 179, 84, 87, 7, 36, 97, 61}
	var contract = Address{110, 77, 84, 109, 194, 103, 190, 117, 2, 7, 51, 62, 159, 37, 108, 17, 228, 174, 205, 167}
	fmt.Println(contract.Hex())
}
