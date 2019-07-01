package staking

import (
	"crypto/sha256"
	"math"
	"math/big"
	"math/rand"

	"github.com/pkg/errors"

	"github.com/BerithFoundation/berith-chain/common"
)

var (
	DIF_MAX = int64(500000)
	DIF_MIN = int64(10000)
	START_IDX = 0
)

const GROUP_UNIT = 100

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
	number     uint64
	period     uint64
	//selections map[uint64]Candidate
	selections []Candidate
}

func NewCandidates(number uint64, period uint64) *Candidates {
	return &Candidates{
		number:     number,
		period:     period,
		selections: make([]Candidate, 0),
	}
}

func (cs *Candidates) Add(c Candidate) {
	cs.selections = append(cs.selections, c)
}

//총 스테이킹 량 , 가산점 추가된 결과
func (cs *Candidates) TotalStakeBalance() *big.Int {
	total := big.NewInt(0)
	for _, c := range cs.selections {
		//adv 적용
		adv := int64(c.GetAdvantage(cs.number, cs.period)*10) + 10
		advStake := new(big.Int).Div(new(big.Int).Mul(c.stake, big.NewInt(adv)), big.NewInt(10))
		total.Add(total, advStake)
	}
	return total.Div(total, big.NewInt(1e+10))
}

//ROI 산출
func (cs *Candidates) GetROI(address common.Address) float64 {
	total := cs.TotalStakeBalance()
	roi := big.NewFloat(0)
	for _, c := range cs.selections {
		if c.address == address {
			stake :=  new(big.Int).Div(c.stake, big.NewInt(1e+8))
			roi.Quo( new(big.Float).SetInt(stake), new(big.Float).SetInt(total))
			break
		}
	}

	f, _ := roi.Float64()
	r := math.Round(f * float64(100)) / float64(100)
	return r
}

//숫자 > 해시 > 숫자
func (cs Candidates) GetSeed(number uint64) int64 {

	bt := []byte{byte(number)}
	hash := sha256.New()
	hash.Write(bt)
	md := hash.Sum(nil)
	h := common.BytesToHash(md)
	seed := h.Big().Int64()

	return seed
}

//BC 선출
func (cs *Candidates) GetBlockCreator(number uint64) *map[common.Address]*big.Int {

	bc := make(map[common.Address]*big.Int, 0)


	//Copy
	var cp *Candidates
	var err error
	size := len(cs.selections)
	if len(cs.selections) < GROUP_UNIT {
		err, cp = nextCandidate(cs, 0, size)
		START_IDX = size
		if err != nil {
			return &bc
		}
	} else {
		err, cp = nextCandidate(cs, 0, GROUP_UNIT)
		START_IDX = GROUP_UNIT
		if err != nil {
			return &bc
		}
	}



	//cp.selections = append(cp.selections, cs.selections...)

	DIF := DIF_MAX
	DIF_R := (DIF_MAX - DIF_MIN) / int64(len(cs.selections))

	seed := cs.GetSeed(number)
	//seed := 100000000 + int64(number)
	//fmt.Println("SEED :: ", seed)
	rand.Seed(seed)

	total := cp.TotalStakeBalance()

	selector := func(value int64) (error, int64, common.Address) {
		temp := new(big.Int).Add(total, big.NewInt(0))
		// Range 확인

		for i:=0; i<len(cp.selections); i++ {
			s := cp.selections[uint64(i)]
			stake := new(big.Int).Div(s.stake, big.NewInt(1e+10))
			temp.Sub(temp, stake)
			if temp.Cmp(big.NewInt(value)) == 1 { //total - stake > value
				continue
			}

			//total - stake <= value
			return nil, int64(i), s.address
		}
		return errors.New("empty SRT"), -1, common.Address{}
	}


	loop := func(value int64) {
		err, key, addr := selector(value)
		if err != nil {
			return
		}

		if DIF == DIF_MAX {
			bc[addr] = big.NewInt(DIF_MAX)
			DIF -= DIF_R
		} else {
			bc[addr] = big.NewInt(DIF)
			DIF -= DIF_R
		}

		stake := new(big.Int).Div(cp.selections[uint64(key)].stake, big.NewInt(1e+10))
		total.Sub(total, stake)

		//remove Index
		//delete(cp.selections, uint64(key))
		cp.selections = removeSlice(cp.selections, key)
	}

	value := int64(0)
	//fmt.Println("TOTAL ::::: ", total.String(), total.Int64())
	for {

		if len(cp.selections) == 0 {
			remainder := len(cs.selections) - START_IDX
			if remainder <= 0 {
				break
			}

			if remainder < GROUP_UNIT {
				err, cp = nextCandidate(cs, START_IDX, START_IDX + remainder)
				START_IDX += remainder
			} else {
				//fmt.Println("NEXT :: ", NEXT)
				err, cp = nextCandidate(cs, START_IDX, START_IDX + GROUP_UNIT)
				START_IDX += GROUP_UNIT
				if err != nil {
					break
				}
			}
			total = cp.TotalStakeBalance()
		}

		if total.Cmp(big.NewInt(0)) == 0 {
			break
		}

		value = rand.Int63n(total.Int64())
		loop(value)
	}
	return &bc
}


func nextCandidate(cs *Candidates, start, end int) (error, *Candidates){
	cp := NewCandidates(cs.number, cs.period)

	nextSize := len(cs.selections[start:end])

	if nextSize == 0 {
		return errors.New("SIZE ZERO"), nil
	}

	if end - start < GROUP_UNIT{
		cp.selections = append(cp.selections, cs.selections[start:]...)
	} else {
		cp.selections = append(cp.selections, cs.selections[start:end]...)
	}
	return nil, cp
}

func removeSlice(cs []Candidate, i int64) []Candidate{
	t1 := cs[:i]
	t2 := cs[i+1:]
	t1 = append(t1, t2...)
	return t1
}
