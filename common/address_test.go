package common

import (
	"fmt"
	"testing"
)

func TestAddressHex(t *testing.T) {
	var coinbase = Address{88, 179, 149, 164, 186, 25, 94, 3, 68, 90, 6, 92, 182, 180, 241, 149, 87, 45, 251, 231}
	fmt.Println(coinbase.Hex())
}
