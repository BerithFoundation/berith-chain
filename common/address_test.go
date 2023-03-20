package common

import (
	"fmt"
	"testing"
)

func TestAddressHex(t *testing.T) {
	// var coinbase = Address{88, 179, 149, 164, 186, 25, 94, 3, 68, 90, 6, 92, 182, 180, 241, 149, 87, 45, 251, 231}
	var msgSender = Address{132, 156, 171, 226, 156, 54, 147, 100, 193, 205, 23, 235, 129, 172, 66, 24, 57, 36, 232, 238}
	// var contract = Address{110, 77, 84, 109, 194, 103, 190, 117, 2, 7, 51, 62, 159, 37, 108, 17, 228, 174, 205, 167}

	// var txHash = Bytes2Hex([]byte{212, 80, 109, 215, 136, 158, 85, 172, 53, 238, 90, 69, 204, 11, 4, 230, 88, 22, 92, 91, 39, 230, 11, 216, 5, 27, 128, 186, 28, 216, 95, 239})
	fmt.Println(msgSender.Hex())
	// fmt.Println(contract.Hex())
}
