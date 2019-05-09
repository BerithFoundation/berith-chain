package staking

import (
	"crypto/sha256"
	"math"
	"math/big"
	"math/rand"

	"github.com/BerithFoundation/berith-chain/common"
)

type Vote struct {
	address common.Address //address
	stake   *big.Int       //stake balance
	block   *big.Int       //block number
	reward  *big.Int       //reward balance
}

func (v *Vote) GetStake() *big.Int {
	return v.stake
}

func (v *Vote) GetReward() *big.Int {
	return v.reward
}

func (v *Vote) GetAdvantage(number, snumber float64, period uint64) float64 {
	//div := 1.2 * (10 ^ 6)
	p := float64(30) / float64(period)
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
func CalcS(votes *[]Vote, number, period uint64) *big.Float {
	stotal := big.NewFloat(0)
	for _, vote := range *votes {
		stake := vote.GetStake()
		reward := vote.GetReward()
		adv := vote.GetAdvantage(float64(number), vote.GetBlockNumber(), period)

		//freward, _ := new(big.Float).Mul(new(big.Float).SetInt(reward), big.NewFloat(0.5)).Int64()
		//s1 := new(big.Int).Add(stake, big.NewInt(freward))
		//s2 := new(big.Int).Mul(s1, big.NewInt(int64(1) + int64(adv)))
		freward := new(big.Float).Mul(new(big.Float).SetInt(reward), big.NewFloat(0.5)) //reward * 0.5
		s1 := new(big.Float).Add(new(big.Float).SetInt(stake), freward)                 //(stake + (reward * 0.5))
		s2 := new(big.Float).Mul(s1, big.NewFloat(1+adv))                               //(stake + (reward * 0.5)) * (1 + adv)

		stotal = new(big.Float).Add(stotal, s2)
	}
	return stotal
}

func CalcP2(votes *[]Vote, stotal *big.Float, number, period uint64) *map[common.Address]int {
	length := len(*votes)

	p := make(map[common.Address]int, length)

	// fmt.Println("******************************LIST & P*********************************")
	for _, vote := range *votes {
		stake := vote.GetStake()
		reward := vote.GetReward()
		adv := vote.GetAdvantage(float64(number), vote.GetBlockNumber(), period)

		//s := (stake + (reward * 0.5)) * (1 + adv)
		freward := new(big.Float).Mul(new(big.Float).SetInt(reward), big.NewFloat(0.5)) //reward * 0.5
		s1 := new(big.Float).Add(new(big.Float).SetInt(stake), freward)                 //(stake + (reward * 0.5))
		s := new(big.Float).Mul(s1, big.NewFloat(1+adv))                                //(stake + (reward * 0.5)) * (1 + adv)

		//temp := s / stotal * 10000000
		temp := new(big.Float).Mul(new(big.Float).Quo(s, stotal), big.NewFloat(10000000))

		tt, _ := temp.Int64()
		if big.NewInt(tt) == big.NewInt(10000000) {
			tt = big.NewInt(9999999).Int64()
		}

		p[vote.address] = int(tt)

		//fmt.Println("\t [BlockNumber]", number)
		//fmt.Print("\t [SIG] : ", vote.address.Hex())
		//fmt.Print("\t [REWARD] : ", reward)
		//fmt.Print("\t [FREWARD] : ", freward)
		//fmt.Print("\t [STAKE] : ", stake)
		//fmt.Print("\t [STOTAL] : ", stotal)
		//fmt.Print("\t [S] : ", s)
		//fmt.Println("\t [P] : ", p[vote.address])
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

		seed := rand.Int63n(9999999)
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
