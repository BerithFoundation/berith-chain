package staking

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"github.com/BerithFoundation/berith-chain/common"
)

func TestVoting2(t *testing.T) {
	number := uint64(10000)
	//epoch := uint64(20)
	perioid := uint64(10)
	loop := 10

	for i:=0; i<10; i++{

		fmt.Println("블록 넘버 :: ", i)

		cs := NewCandidates(number, perioid)

		v := rand.Int63n(100000000000000000)
		fmt.Println(v)



		for i := 0; i < loop; i++ {

			value := int64(100000)

			//stake := new(big.Int).Mul(big.NewInt(10000000 + (int64(i) * 1)), big.NewInt(1e+18))
			stake := new(big.Int).Mul(big.NewInt(value + int64(i)), big.NewInt(1e+18))
			c := Candidate{common.BytesToAddress([]byte(strconv.Itoa(i))), stake, big.NewInt(1), big.NewInt(100)}
			cs.Add(c)
		}


		bc := cs.GetBlockCreator(number + uint64(i))
		//bc := cs.GetBlockCreator(number + uint64(i), epoch, perioid)

		for key, val := range *bc {
			fmt.Print("SIGNER :: ", common.Bytes2Hex(key.Bytes()))
			fmt.Println(" VALUE :: ", val)
		}
		fmt.Println("=================================================================================")
	}




}