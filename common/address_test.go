package common

import (
	"fmt"
	"testing"
)

func TestAddressHex(t *testing.T) {
	// var coinbase = Address{88, 179, 149, 164, 186, 25, 94, 3, 68, 90, 6, 92, 182, 180, 241, 149, 87, 45, 251, 231}
	// var contract = Address{110, 77, 84, 109, 194, 103, 190, 117, 2, 7, 51, 62, 159, 37, 108, 17, 228, 174, 205, 167}

	var msgSender = Address{35, 69, 191, 119, 209, 222, 158, 172, 246, 111, 232, 26, 9, 168, 108, 250, 178, 18, 165, 66}
	fmt.Println("Sender : ", msgSender.Hex())
	var txHash = Bytes2Hex([]byte{201, 126, 158, 129, 225, 64, 119, 91, 182, 0, 234, 225, 149, 47, 211, 88, 92, 192, 232, 21, 140, 180, 68, 217, 94, 96, 103, 0, 239, 122, 111, 91})
	fmt.Println("Tx : ", txHash)
}
