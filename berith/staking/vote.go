package staking

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"math/big"
	"math/rand"

	"github.com/BerithFoundation/berith-chain/common"
)

type Candidate struct {
	address common.Address //address
	stake   *big.Int       //stake balance
	block   *big.Int       //block number -- Contribution
	reward  *big.Int       //reward balance
}


func (c *Candidate) GetStake() *big.Int {
	return c.stake
}

func (c *Candidate) GetReward() *big.Int {
	return c.reward
}

func (c *Candidate) GetBlockNumber() float64 {
	return float64(c.block.Uint64())
}

//Stake 기간 Adv를 구한다.
func (c *Candidate) GetAdvantage(number uint64, period uint64) float64 {
	p := float64(30) / float64(period) //30초 기준의 공식이기때문에
	y := 1.2 * float64(p)
	div := y * math.Pow(10, 6) //10의6승

	adv := (float64(number) - float64(c.block.Uint64())) / div
	if adv >= 1 {
		return 1
	} else {
		return adv
	}
}


///////////////////////////////////////////////////////////////////////////////////////////
type Candidates struct {
	number uint64
	period uint64
	selections map[uint64]Candidate
}

func NewCandidates(number uint64, period uint64) *Candidates {
	return &Candidates{
		number: number,
		period: period,
		selections: make(map[uint64]Candidate, 0),
	}
}

func (cs *Candidates) Add(c Candidate) {
	s := len(cs.selections)
	cs.selections[uint64(s)] = c
}

func (cs *Candidates) Remove(key uint64){
	delete(cs.selections, key)
}

//총 스테이킹 량 , 가산점 추가된 결과
func (cs *Candidates) TotalStakeBalance() *big.Int {
	total := big.NewInt(0)
	for _,c := range cs.selections {
		//adv 적용
		adv := int64(c.GetAdvantage(cs.number, cs.period) * 10) + 10
		advStake := new(big.Int).Div(new(big.Int).Mul(c.stake, big.NewInt(adv)), big.NewInt(10))
		total = new(big.Int).Add(total, advStake)
	}
	return new(big.Int).Div(total, big.NewInt(1e+10))
}

type StakerRange struct {
	address common.Address
	val *big.Int
}
//Make Staker Range Table
func (cs *Candidates)MakeSRT() (*big.Int, *map[common.Address]StakerRange) {
	srt := make(map[common.Address]StakerRange, 0)
	total := big.NewInt(0)

	for _, c := range cs.selections{
		//ADV
		adv := int64(c.GetAdvantage(cs.number, cs.period) * 10) + 10
		advStake := new(big.Int).Div(new(big.Int).Mul(c.stake, big.NewInt(adv)), big.NewInt(10))
		total = new(big.Int).Add(total, advStake)

		sr := StakerRange{
			address: c.address,
			val: new(big.Int).Div(advStake, big.NewInt(1e+10)),
		}

		srt[c.address] = sr
	}

	return new(big.Int).Div(total, big.NewInt(1e+10)), &srt

}


//BC 선출
func (cs *Candidates)GetBlockCreator(number uint64, epoch, period uint64) *map[common.Address]*big.Int {

	bc := make(map[common.Address]*big.Int, 0)

	seed := 100000000 + int64(number)
	fmt.Println("SEED :: ", seed)
	rand.Seed(seed)

	total, srt := cs.MakeSRT() //SRT 만들기

	selector :=  func (value int64, srt *map[common.Address]StakerRange) (error, *StakerRange){
		r := big.NewInt(0)
		// Range 확인
		for _, s := range *srt{
			r = new(big.Int).Add( r, s.val) //range

			if r.Cmp(big.NewInt(value)) != 1{
				continue
			}

			total = new(big.Int).Sub(total, s.val)
			delete(*srt, s.address)

			return nil, &s
		}
		return errors.New("empty SRT"), nil
	}


	for {
		if len(*srt) == 0{
			break
		}

		value := rand.Int63n(total.Int64())

		_, sr := selector(value, srt)

		bc[sr.address] = big.NewInt(1000)


	}


	return &bc
}

//func (cs *Candidates)ranInsert(r, adv float64, bc *map[common.Address]*big.Int ) bool {
//   value := rand.Int63n(int64(r))
//   //fmt.Print("RANDOM IDX :: ", idx)
//   //fmt.Println("RANDOM VALUE :: ", value)
//
//   for i:=0; i<len(cs.selections); i++{
//
//      c := cs.selections[uint64(i)]
//
//
//      t1 := int64(float64(10) + (adv * float64(10))) // 1 + adv
//      t2 := new(big.Int).Mul(c.stake, big.NewInt(t1))
//      sv := new(big.Int).Div(new(big.Int).Div(t2, big.NewInt(1e+18)), big.NewInt(10))
//   }
//
//   for i, c := range cs.selections {
//      //Staking Balance 가산점 적용
//      t1 := int64(float64(10) + (adv * float64(10))) // 1 + adv
//      t2 := new(big.Int).Mul(c.stake, big.NewInt(t1))
//      sv := new(big.Int).Div(new(big.Int).Div(t2, big.NewInt(1e+18)), big.NewInt(10))
//
//      if value < c.max
//
//      //테이블 셀렉트
//      if value < sv.Int64() {
//         //DIF 도 여기에 Insert
//         if _, exists := (*bc)[c.address]; !exists {
//            (*bc)[c.address] = sv
//            return true
//         }
//      }
//
//   }
//
//   return false
//}


//S구하기
//func CalcS(votes *[]Vote, number, period uint64) *big.Float {
//   stotal := big.NewFloat(0)
//   for _, vote := range *votes {
//      stake := vote.GetStake()
//      reward := big.NewInt(0)//vote.GetReward()
//      //adv := vote.GetAdvantage(float64(number), vote.GetBlockNumber(), period)
//
//      //freward, _ := new(big.Float).Mul(new(big.Float).SetInt(reward), big.NewFloat(0.5)).Int64()
//      //s1 := new(big.Int).Add(stake, big.NewInt(freward))
//      //s2 := new(big.Int).Mul(s1, big.NewInt(int64(1) + int64(adv)))
//      freward := new(big.Float).Mul(new(big.Float).SetInt(reward), big.NewFloat(0.5)) //reward * 0.5
//      s1 := new(big.Float).Add(new(big.Float).SetInt(stake), freward) //(stake + (reward * 0.5))
//      s2 := new(big.Float).Mul(s1, big.NewFloat(1 + 0)) //(stake + (reward * 0.5)) * (1 + adv)
//
//      stotal = new(big.Float).Add(stotal, s2)
//   }
//   return stotal
//}

//
//func CalcP2(votes *[]Vote, stotal *big.Float, number, period uint64) *map[common.Address]int {
//   length := len(*votes)
//
//   p := make(map[common.Address]int, length)
//
//   // fmt.Println("******************************LIST & P*********************************")
//   for _, vote := range *votes {
//      stake := vote.GetStake()
//      reward := big.NewInt(0)//vote.GetReward()
//      adv := vote.GetAdvantage(float64(number), vote.GetBlockNumber(), period)
//
//      //s := (stake + (reward * 0.5)) * (1 + adv)
//      freward := new(big.Float).Mul(new(big.Float).SetInt(reward), big.NewFloat(0.5)) //reward * 0.5
//      s1 := new(big.Float).Add(new(big.Float).SetInt(stake), freward) //(stake + (reward * 0.5))
//      s := new(big.Float).Mul(s1, big.NewFloat(1 + 0)) //(stake + (reward * 0.5)) * (1 + adv)
//
//      //temp := s / stotal * 10000000
//      temp := new(big.Float).Mul(new(big.Float).Quo(s, stotal),  big.NewFloat(10000000))
//
//      tt, _ := temp.Int64()
//      if big.NewInt(tt) == big.NewInt(10000000) {
//         tt = big.NewInt(9999999).Int64()
//      }
//
//      p[vote.address] = int(tt)
//
//      //fmt.Println("\t [BlockNumber]", number)
//      //fmt.Print("\t [SIG] : ", vote.address.Hex())
//      //fmt.Print("\t [REWARD] : ", reward)
//      //fmt.Print("\t [FREWARD] : ", freward)
//      //fmt.Print("\t [STAKE] : ", stake)
//      //fmt.Print("\t [STOTAL] : ", stotal)
//      //fmt.Print("\t [S] : ", s)
//      //fmt.Println("\t [P] : ", p[vote.address])
//   }
//
//   // fmt.Println("***********************************************************************")
//
//   return &p
//}
//
//func CalcR2(votes *[]Vote, p *map[common.Address]int) *[]int {
//   length := len(*votes)
//   r := make([]int, 0)
//   for i := 0; i < length; i++ {
//      r = append(r, 0)
//      for j := 0; j <= i; j++ {
//         addr := (*votes)[j].address
//         r[i] += (*p)[addr]
//      }
//   }
//   return &r
//}
//
//func GetSigners(seed int64, votes *[]Vote, r *[]int, epoch uint64) *[]common.Address {
//   sigs := make([]common.Address, 0)
//   for i := 0; uint64(i) < epoch; i++ {
//      if i == 0 {
//         rand.Seed(seed + int64(i))
//      } else {
//         a := []byte{byte(seed + int64(i))}
//         //sum := sha256.Sum256(a)
//         hash := sha256.New()
//         hash.Write(a)
//         md := hash.Sum(nil)
//
//         h := common.BytesToHash(md)
//         //mdStr := hex.EncodeToString(md)
//         newSeed := common.HexToAddress(h.Hex()).Big().Int64()
//         rand.Seed(newSeed)
//      }
//
//      seed := rand.Int63n(9999999)
//      //seed := int64(876543)
//      //fmt.Println("SEED", seed)
//      for i, v := range *votes {
//         if seed < int64((*r)[i]) {
//            //fmt.Printf("%s \n",v.address)
//            sigs = append(sigs, v.address)
//            break
//         }
//      }
//   }
//
//   return &sigs
//}
