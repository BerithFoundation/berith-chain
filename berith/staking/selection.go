/**
[BERITH]
선출 연산을 담당 하는 go파일


 */

package staking

import (
	"crypto/sha256"
	"errors"
	"math"
	"math/big"
	"math/rand"

	"github.com/BerithFoundation/berith-chain/common"
)

const (
	MAX_MINERS = 22
)

var (
	DIF_MAX = int64(5000000)
	DIF_MIN = int64(10000)
)

/**
[BERITH]
선출을 위해 Staking 한 계정들의 정보를 담는 구조체
 */
type Candidate struct {
	address  common.Address //계정 주소
	stake    uint64         //stake balance
	block    uint64         //block number -- Contribution (staking 한 블록 넘버)
	reward   uint64         //reward balance
	val      uint64         //total staking balance
	advStake uint64         //advantage staking balance
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
/**
[BERITH]

 */
type Candidates struct {
	number uint64
	period uint64
	selections []Candidate
	total      uint64 //Total Staking  + Adv
	ts         uint64 //Total Staking Value
}

func NewCandidates(number uint64, period uint64) *Candidates {
	return &Candidates{
		number:     number,
		period:     period,
		selections: make([]Candidate, 0),
		total:      0,
		ts:         0,
	}
}

/*
[BERITH]
BC 선출을 하기 위해 Staker 를 등록하기 위한 함수
이후에 호출 될 함수는 BlockCreator 함수이다.
 */
func (cs *Candidates) Add(c Candidate) {
	adv := uint64(c.GetAdvantage(cs.number, cs.period)*10) + 10
	c.advStake = c.stake * adv
	cs.total += c.advStake
	c.val = cs.total
	cs.selections = append(cs.selections, c)

	cs.ts += c.stake //Total Staking
}

/*
[BERITH]
블록 넘버를 해시로 바꾸고 그것을 강제로 int64로 변경 하는 함수
결과 값을 Seed 로 쓴다.
 */
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
type VoteResult struct {
	Score *big.Int `json:"score"`
	Rank  int      `json:"rank"`
}

/*
[BERITH]
랜덤값으로 binarySearch 하기 위한 원형큐 구조체
*/
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

/**
[BERITH]
Random 값을 폭단위로 binarySearch 한다.
 */
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

/*
[BERITH]
BC 선출을 하기 위한 함수
선출된 BC map 을 리턴 한다.
 */
func (cs *Candidates) BlockCreator(number uint64) *map[common.Address]VoteResult {
	queue := &Queue{
		storage: make([]Range, len(cs.selections)),
		size:    len(cs.selections) + 1,
		front:   0,
		rear:    0,
	}
	result := make(map[common.Address]VoteResult)

	DIF := DIF_MAX
	DIF_R := (DIF_MAX - DIF_MIN) / int64(len(cs.selections))

	rand.Seed(cs.GetSeed(number))

	queue.enqueue(Range{
		min:   0,
		max:   cs.total,
		start: 0,
		end:   len(cs.selections),
	})
	for count := 1; count <= MAX_MINERS && queue.front != queue.rear; count++ {
		r, _ := queue.dequeue()
		account := r.binarySearch(queue, cs)
		result[account] = VoteResult{
			Score: big.NewInt(DIF + int64(cs.ts)),
			Rank:  count,
		}
		DIF -= DIF_R

	}

	//fmt.Println(DIF)

	return &result
}

/*
[BERITH]
예상 BC 선출 비율을 계산하는 함수
*/
func (cs *Candidates) getJoinRatio(address common.Address) float64 {
	stake := uint64(0)
	for _, c := range cs.selections {
		if c.address == address {
			stake = c.advStake
			break
		}
	}

	f := float64(stake) / float64(cs.total)
	r := math.Round(f * float64(100))
	return r
}
