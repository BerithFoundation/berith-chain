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
	DIF_MAX   = int64(500000)
	DIF_MIN   = int64(10000)
	TS = uint64(0) //Total Staking Value
)

type Candidate struct {
	address common.Address //address
	stake   uint64         //stake balance
	block   uint64         //block number -- Contribution
	reward  uint64         //reward balance
	val     uint64		   //sum
	advStake		uint64		   //advStake
}

func (c *Candidate) GetStake() uint64 {
	return c.stake
}

func (c *Candidate) GetReward() uint64 {
	return c.reward
}

func (c *Candidate) GetBlockNumber() float64 {
	return float64(c.block)
}

//Stake 기간 Adv를 구한다.
func (c *Candidate) GetAdvantage(number uint64, period uint64) float64 {
	p := float64(30) / float64(period) //30초 기준의 공식이기때문에
	y := 1.2 * float64(p)
	div := y * math.Pow(10, 6) //10의6승

	adv := (float64(number) - c.GetBlockNumber()) / div
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
	//selections map[uint64]Candidate
	selections []Candidate
	total      uint64
}

func NewCandidates(number uint64, period uint64) *Candidates {
	return &Candidates{
		number:     number,
		period:     period,
		selections: make([]Candidate, 0),
		total:      0,
	}
}

func (cs *Candidates) Add(c Candidate) {
	adv := uint64(c.GetAdvantage(cs.number, cs.period) * 10) + 10
	c.advStake = c.stake * adv
	cs.total += c.advStake
	c.val = cs.total
	cs.selections = append(cs.selections, c)

	TS += c.stake //Total Staking 
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

type Range struct {
	min   uint64
	max   uint64
	start int
	end   int
}

type Queue struct {
	storage []Range
	size    int
	front   int
	rear    int
}

func (q *Queue) enqueue(r Range) error {
	next := (q.rear + 1) % q.size
	if next == q.front {
		return errors.New("Queue is full")
	}
	q.storage[q.rear] = r
	q.rear = next
	return nil
}

func (q *Queue) dequeue() (Range, error) {
	if q.front == q.rear {
		return Range{}, errors.New("Queue is Empty")
	}
	result := q.storage[q.front]
	q.front = (q.front + 1) % q.size
	return result, nil
}

func (r Range) binarySearch(q *Queue, cs *Candidates) common.Address {
	if r.end-r.start <= 1 {
		return cs.selections[r.start].address
	}

	random := uint64(rand.Int63n(int64(r.max-r.min))) + r.min

	start := r.start
	end := r.end
	for {
		target := (start + end) / 2
		a := r.min
		if target > 0 {
			a = cs.selections[target-1].val
		}
		b := cs.selections[target].val

		if random >= a && random <= b {
			if r.start != target {
				q.enqueue(Range{
					min:   r.min,
					max:   a - 1,
					start: r.start,
					end:   target,
				})
			}
			if target+1 != r.end {
				q.enqueue(Range{
					min:   b + 1,
					max:   r.max,
					start: target + 1,
					end:   r.end,
				})
			}
			return cs.selections[target].address
		}

		if random < a {
			end = target
		} else {
			start = target + 1
		}
	}
}

func (cs *Candidates) BlockCreator(number uint64) *map[common.Address]*big.Int {
	queue := &Queue{
		storage: make([]Range, len(cs.selections)),
		size:    len(cs.selections) + 1,
		front:   0,
		rear:    0,
	}
	result := make(map[common.Address]*big.Int)

	DIF := DIF_MAX
	DIF_R := (DIF_MAX - DIF_MIN) / int64(len(cs.selections))

	rand.Seed(cs.GetSeed(number))

	queue.enqueue(Range{
		min:   0,
		max:   cs.total,
		start: 0,
		end:   len(cs.selections),
	})

	for queue.front != queue.rear {
		r, _ := queue.dequeue()
		account := r.binarySearch(queue, cs)
		result[account] = big.NewInt(DIF + int64(TS))
		DIF -= DIF_R
	}

	//fmt.Println(DIF)

	return &result
}


//ROI 산출
func (cs *Candidates) getJoinRatio(address common.Address) float64 {
	stake := uint64(0)
	for _, c := range cs.selections {
		if c.address == address {
			stake =  c.advStake
			break
		}
	}

	f := float64(stake) / float64(cs.total)
	r := math.Round(f * float64(100))
	return r
}