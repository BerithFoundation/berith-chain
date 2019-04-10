package staking

import (
	"crypto/sha256"
	"math"
	"math/big"
	"math/rand"

	"bitbucket.org/ibizsoftware/berith-chain/common"
)

type Vote struct {
	address common.Address //address
	stake   *big.Int       //stake balance
	block   *big.Int       //block number
	reward  *big.Int       //reward balance
}

func (v *Vote) GetStake() float64 {
	return float64(v.stake.Uint64())
}

func (v *Vote) GetReward() float64 {
	return float64(v.reward.Uint64())
}

func (v *Vote) GetAdvantage(number, snumber float64, perioid uint64) float64 {
	//div := 1.2 * (10 ^ 6)
	p := float64(30) / float64(perioid)
	y := 1.2 * float64(p)
	div := y * math.Pow(10, 6)
	adv := (number - snumber) / div
	if adv >= 1 {
		return 1
	} else {
		return adv
	}
}

func (v *Vote) GetBlockNumber() float64 {
	return float64(v.block.Uint64())
}

///////////////////////////////////////////////////////////////////////////////////////////

//S구하기
func CalcS(votes *[]Vote, number, perioid uint64) float64 {
	var stotal float64 = 0
	for _, vote := range *votes {
		stake := vote.GetStake()
		reward := vote.GetReward()
		adv := vote.GetAdvantage(float64(number), vote.GetBlockNumber(), perioid)
		s := (stake + (reward * 0.5)) * (1 + adv)
		stotal += s
	}
	return stotal
}

//P구하기
func CalcP(votes *[]Vote, stotal float64, number, perioid uint64) *[]int {
	length := len(*votes)
	p := make([]int, length)

	for i, vote := range *votes {
		stake := vote.GetStake()
		reward := vote.GetReward()
		adv := vote.GetAdvantage(float64(number), vote.GetBlockNumber(), perioid)
		s := (stake + (reward * 0.5)) * (1 + adv)
		temp := s / stotal * 1000000
		if temp == 1000000 {
			temp = 999999
		}
		p[i] = int(temp)
		// fmt.Println("******************************LIST & P*********************************")
		// fmt.Print("[SIG] : ", vote.address.Hex())
		// fmt.Println("\t [P] : ", p[i])
		// fmt.Println("***********************************************************************")
	}

	return &p
}

func CalcP2(votes *[]Vote, stotal float64, number, perioid uint64) *map[common.Address]int {
	length := len(*votes)

	p := make(map[common.Address]int, length)

	// fmt.Println("******************************LIST & P*********************************")
	for _, vote := range *votes {
		stake := vote.GetStake()
		reward := vote.GetReward()
		adv := vote.GetAdvantage(float64(number), vote.GetBlockNumber(), perioid)
		s := (stake + (reward * 0.5)) * (1 + adv)
		temp := s / stotal * 1000000
		if temp == 1000000 {
			temp = 999999
		}
		p[vote.address] = int(temp)
		// fmt.Print("[SIG] : ", vote.address.Hex())
		// fmt.Println("\t [P] : ", p[vote.address])
	}

	// fmt.Println("***********************************************************************")

	return &p
}

func CalcR2(votes *[]Vote, p *map[common.Address]int) *[]int {
	length := len(*votes)
	r := make([]int, 0)
	for i := 0; i < length; i++ {
		r = append(r, 0)
		for j := 0; j <= i; j++ {
			addr := (*votes)[j].address
			r[i] += (*p)[addr]
		}
	}
	return &r
}

func CalcR(votes *[]Vote, p *[]int) *[]int {
	length := len(*votes)
	r := make([]int, 0)
	for i := 0; i < length; i++ {
		r = append(r, 0)
		for j := 0; j <= i; j++ {
			r[i] += (*p)[j]
		}
	}
	return &r
}

func GetSigners(seed int64, votes *[]Vote, r *[]int, epoch uint64) *[]common.Address {
	sigs := make([]common.Address, 0)
	for i := 0; uint64(i) < epoch; i++ {
		if i == 0 {
			rand.Seed(seed + int64(i))
		} else {
			a := []byte{byte(seed + int64(i))}
			//sum := sha256.Sum256(a)
			hash := sha256.New()
			hash.Write(a)
			md := hash.Sum(nil)

			h := common.BytesToHash(md)
			//mdStr := hex.EncodeToString(md)
			newSeed := common.HexToAddress(h.Hex()).Big().Int64()
			rand.Seed(newSeed)
		}

		seed := rand.Int63n(999999)
		//seed := int64(876543)
		//fmt.Println("SEED", seed)
		for i, v := range *votes {
			if seed < int64((*r)[i]) {
				//fmt.Printf("%s \n",v.address)
				sigs = append(sigs, v.address)
				break
			}
		}
	}

	return &sigs
}
