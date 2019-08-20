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

	number := uint64(10000)
	//epoch := uint64(20)
	perioid := uint64(10)
	loop := 100

	cs := NewCandidates(number, perioid)

	for i := 0; i < loop; i++ {

		value := uint64(10000000)

		//stake := new(big.Int).Mul(big.NewInt(10000000 + (int64(i) * 1)), big.NewInt(1e+18)
		c := Candidate{common.BytesToAddress([]byte(strconv.Itoa(i))), value, 1, 0, 0, 0}
		cs.Add(c)
	}
	bc := cs.BlockCreator(number)
	//cs.GetBlockCreator(number)
	//bc := cs.GetBlockCreator(number)
	//
	fmt.Println(len(*bc))
	//
	idx := 1
	for key, val := range *bc {
		fmt.Print("SIGNER "+strconv.Itoa(idx)+"::  ", common.Bytes2Hex(key.Bytes()))
		fmt.Println(" SCORE :: ", val.score.String())
		fmt.Println(" RANK :: ", val.rank)
		idx++
	}

}

//func TestVoting2(t *testing.T) {
//	number := uint64(10000)
//	epoch := uint64(20)
//	perioid := uint64(10)
//	loop := 10000
//
//	for i:=0; i<10; i++{
//
//		fmt.Println("블록 넘버 :: ", i)
//
//		cs := NewCandidates(number, perioid)
//
//		v := rand.Int63n(100000000000000000)
//		fmt.Println(v)
//
//
//
//		for i := 0; i < loop; i++ {
//
//			value := int64(100000)
//
//			//stake := new(big.Int).Mul(big.NewInt(10000000 + (int64(i) * 1)), big.NewInt(1e+18))
//			stake := new(big.Int).Mul(big.NewInt(value + int64(i)), big.NewInt(1e+18))
//			c := Candidate{common.BytesToAddress([]byte(strconv.Itoa(i))), stake, big.NewInt(1), big.NewInt(100)}
//			cs.Add(c)
//		}
//
//
//		cs.GetBlockCreator(number + uint64(i), epoch, perioid)
//		//bc := cs.GetBlockCreator(number + uint64(i), epoch, perioid)
//
//		//for key, val := range *bc {
//		//	fmt.Print("SIGNER :: ", common.Bytes2Hex(key.Bytes()))
//		//	fmt.Println(" VALUE :: ", val)
//		//}
//		//fmt.Println("=================================================================================")
//	}
//}
