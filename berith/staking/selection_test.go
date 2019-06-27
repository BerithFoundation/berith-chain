package staking

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/BerithFoundation/berith-chain/common"
)

func TestVoting2(t *testing.T) {

	//rand.Seed(100)
	//fmt.Println(rand.Int63n(1000))
	//fmt.Println(rand.Int63n(1000))
	//fmt.Println(rand.Int63n(1000))
	//fmt.Println(rand.Int63n(1000))

	number := uint64(10000)
	//epoch := uint64(20)
	perioid := uint64(10)
	loop := 4

	cs := NewCandidates(number, perioid)

	for i := 0; i < loop; i++ {

		value := int64(100000)

		//stake := new(big.Int).Mul(big.NewInt(10000000 + (int64(i) * 1)), big.NewInt(1e+18))
		stake := new(big.Int).Mul(big.NewInt(value), big.NewInt(1e+18))
		c := Candidate{common.BytesToAddress([]byte(strconv.Itoa(i))), stake, big.NewInt(1), big.NewInt(100)}
		cs.Add(c)
	}


	//cs.GetBlockCreator(number)
	bc := cs.GetBlockCreator(number)

	idx := 1
	for key, val := range *bc {
		fmt.Print("SIGNER " + strconv.Itoa(idx) + "::  ", common.Bytes2Hex(key.Bytes()))
		fmt.Println(" VALUE :: ", val)
		idx++
	}




}
