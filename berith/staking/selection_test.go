package staking

import (
	"fmt"
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

	number := uint64(100000)
	//epoch := uint64(20)
	perioid := uint64(10)
	loop := 6

	cs := NewCandidates(number, perioid)

	for i := 0; i < loop; i++ {

		value := uint64(100000)

		//stake := new(big.Int).Mul(big.NewInt(10000000 + (int64(i) * 1)), big.NewInt(1e+18)
		c := Candidate{common.BytesToAddress([]byte(strconv.Itoa(i))), value, 1, 0, 0}
		cs.Add(c)
	}
	bc := cs.BinarySearch(number)
	//cs.GetBlockCreator(number)
	//bc := cs.GetBlockCreator(number)
	//
	fmt.Println(len(*bc))

	idx := 1
	for key, val := range *bc {
		fmt.Print("SIGNER "+strconv.Itoa(idx)+"::  ", common.Bytes2Hex(key.Bytes()))
		fmt.Println(" VALUE :: ", val)
		idx++
	}

}
