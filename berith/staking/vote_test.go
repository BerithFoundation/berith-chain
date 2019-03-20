package staking

import (
	"bitbucket.org/ibizsoftware/berith-chain/common"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type VoteTest struct {
	address common.Address //address
	stake *big.Int	//stake balance
	block *big.Int //stake block number
	reward *big.Int //reward balance
}

func TestVoting(t *testing.T){

	number := 100

	temps := make([]VoteTest, 0)

	for i:=0; i<1360; i++ {
		v := VoteTest{common.BytesToAddress([]byte(strconv.Itoa(i))), big.NewInt(10000000 + int64(i)), big.NewInt(1), big.NewInt(100)}
		temps = append(temps, v)
	}

	length := len(temps)

	//S 계산
	var stotal float64 = 0
	for _, vote := range temps {
		s := (float64(vote.stake.Uint64()) + float64(vote.reward.Uint64())*0.5) * (1 + fadvTest(float64(number), float64(vote.block.Uint64())))
		stotal += s
	}
	//fmt.Println("S :: ", stotal)


	//P 계산
	p := make([]int, length)
	for i, vote := range temps {
		s := (float64(vote.stake.Uint64()) + float64(vote.reward.Uint64())*0.5) * (1 + fadvTest(float64(number), float64(vote.block.Uint64())))
		temp:= s/ stotal * 1000000
		if temp == 1000000 {
			temp = 999999
		}

		p[i] = int(temp)

		//p = append(p, int(temp))
	}
	fmt.Println("P :: ", p)

	//R 계산
	//r := make([]int, length)
	//for i, _ := range temps {
	//	for j:=0; j<=i; j++ {
	//
	//		r[i] += p[j]
	//	}
	//}

	//r := make([]int, length)
	//for i, _ := range temps {
	//	for j:=0; j<=i; j++ {
	//
	//		r[i] += p[j]
	//	}
	//}

	r := make([]int, 0)
	for i:=0; i<length; i++{
		r = append(r, 0)
		for j:=0; j<=i; j++ {
			r[i] += p[j]
		}
	}

	//fmt.Println("R :: ", r)


	n := common.HexToAddress("0x2c21bf2f10eb55d538f1af154260025f605613283437d872f9ede4736b41a58d").Big().Int64()
	//fmt.Println(n)

	signers := make([]common.Address, 0)

	for i:=0; i<360; i++ {
		rand.Seed(n + int64(i))
		seed := rand.Int63n(999999)
		//seed := int64(876543)
		//fmt.Println("SEED", seed)
		for i, v := range temps {
			if seed < int64(r[i]) {
				//fmt.Printf("%s \n",v.address)
				signers = append(signers, v.address)
				break
			}
		}
	}


	//for _, sig := range signers {
	//	fmt.Println("SIGNER :: ", common.Bytes2Hex(sig.Bytes()))
	//}
	//fmt.Println("SIGNER :: ", signers)





	//for _, vote := range temps{
	//	s:= (float64(vote.stake.Uint64()) + float64(vote.reward.Uint64())*0.5) * (1 + fadv(float64(number), float64(vote.block.Uint64())))
	//
	//
	//	p := new(big.Int).Div(total, vote.stake)
	//
	//	fmt.Printf("%s [보정수치] %f , [P] %d \r\n", vote.address, s, p.Uint64())
	//
	//	r := 0
	//	for i:=1;  uint64(i) < p.Uint64(); i++ {
	//		r = i + r
	//	}
	//
	//	fmt.Println("R :: ", r)
	//
	//}


}


//ADV
func fadvTest(number, snumber float64) float64 {

	div := 1.2 * math.Pow(10, 6)

	adv := (number - snumber) / div
	if adv >= 1 {
		return 1
	} else {
		return adv
	}
}

func TestVoting2(t *testing.T)  {
	number := uint64(100)

	votes := make([]Vote, 0)

	for i:=0; i<1360; i++ {
		v := Vote{common.BytesToAddress([]byte(strconv.Itoa(i))), big.NewInt(10000000), big.NewInt(1), big.NewInt(100)}
		votes = append(votes, v)
	}

	stotal := CalcS(&votes, number)
	p := CalcP(&votes, stotal, number)
	fmt.Println(*p)
	r := CalcR(&votes, p)

	n := common.HexToAddress("0x2c21bf2f10eb55d538f1af154260025f605613283437d872f9ede4736b41a58d").Big().Int64()

	GetSigners(n, &votes, r, 20)

	//for _, sig := range *signers {
	//	fmt.Println("SIGNER :: ", common.Bytes2Hex(sig.Bytes()))
	//}
}

func TestReward(t *testing.T){
	for i:=0; i<100000000; i+=10000 {
		r := reward(uint64(i))

		temp := r * 1e+10
		re := new(big.Int).Mul(big.NewInt(int64(temp)), big.NewInt(1e+8))
		//re := big.NewInt(int64(temp))

		fmt.Println(re)

		time.Sleep(100)
	}

}

func reward(number uint64) float64 {
	up := 5.5 * 100 * math.Pow(10, 7.2)
	down := float64(number) + math.Pow(10, 7.8)

	y := up/down - 60.0

	if y < 0 {
		return float64(0)
	}
	return y
}