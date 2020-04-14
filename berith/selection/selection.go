/**
[BERITH]
선출 연산을 담당 하는 go파일
*/
package selection

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"sort"

	"github.com/BerithFoundation/berith-chain/params"

	"github.com/BerithFoundation/berith-chain/berith/staking"
	"github.com/BerithFoundation/berith-chain/core/state"

	"github.com/BerithFoundation/berith-chain/common"
)

const (
	MAX_MINERS = 10000
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
	address common.Address //계정 주소
	point   uint64         //계정의 포인트 (내가 뽑힐 확률 : 나의 포인트 / 전체 유저의 포인트)
	val     uint64         //블록 생성자 선출을 위해 사용되는 값
}

func (c *Candidate) GetPoint() uint64 {
	return c.point
}

///////////////////////////////////////////////////////////////////////////////////////////
/**
[BERITH]

*/
type Candidates struct {
	selections []Candidate
	total      uint64 //Total Staking  + Adv
	ts         uint64
}

func NewCandidates() *Candidates {
	return &Candidates{
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
	cs.total += c.point
	c.val = cs.total
	cs.selections = append(cs.selections, c)
}

/*
[BERITH]
블록 넘버를 해시로 바꾸고 그것을 강제로 int64로 변경 하는 함수
결과 값을 Seed 로 쓴다.
*/
func (cs Candidates) GetSeed(config *params.ChainConfig, number uint64) int64 {

	bt := []byte{byte(number)}
	if config.IsBIP2(big.NewInt(0).SetUint64(number)) {
		bt = big.NewInt(0).SetUint64(number).Bytes()
	}
	hash := sha256.New()
	hash.Write(bt)
	md := hash.Sum(nil)
	h := common.BytesToHash(md)
	seed := h.Big().Int64()

	return seed
}

func (cs Candidates) GetBIP2Seed(number uint64) int64 {

	bt := big.NewInt(int64(number)).Bytes()
	hash := sha256.New()
	hash.Write(bt)
	md := hash.Sum(nil)
	h := common.BytesToHash(md)
	seed := h.Big().Int64()

	return seed
}

func (cs *Candidates) selectBlockCreator(config *params.ChainConfig, number uint64) VoteResults {

	queue := &Queue{
		storage: make([]Range, len(cs.selections)),
		size:    len(cs.selections) + 1,
		front:   0,
		rear:    0,
	}
	result := make(VoteResults)

	DIF := DIF_MAX
	DIF_R := (DIF_MAX - DIF_MIN) / int64(len(cs.selections))

	rand.Seed(cs.GetSeed(config, number))

	_ = queue.enqueue(Range{
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
	return result
}

func (cs *Candidates) selectBIP3BlockCreator(config *params.ChainConfig, number uint64) VoteResults {

	result := make(VoteResults)

	DIF := DIF_MAX
	DIF_R := (DIF_MAX - DIF_MIN) / int64(len(cs.selections))
	rank := 1
	rand.Seed(cs.GetSeed(config, number))

	for len(cs.selections) > 0 {

		target := uint64(rand.Int63n(int64(cs.total)))

		var chosen int
		start := 0
		end := len(cs.selections) - 1

		for {
			mid := (start + end) / 2
			a := uint64(0)
			if mid > 0 {
				a = cs.selections[mid-1].val
			}
			b := cs.selections[mid].val

			if target >= a && target <= b {
				chosen = mid
				cddt := cs.selections[mid]
				result[cddt.address] = VoteResult{
					Rank:  rank,
					Score: big.NewInt(DIF),
				}
				DIF -= DIF_R
				rank++
				break
			}

			if target < a {
				end = mid - 1
			}
			if target > b {
				start = mid + 1
			}
		}

		out := cs.selections[chosen]
		for i := chosen; i+1 < len(cs.selections); i++ {
			newCddt := cs.selections[i+1]
			newCddt.val -= out.point
			cs.selections[i] = newCddt
		}

		cs.selections = cs.selections[:len(cs.selections)-1]
		cs.total -= out.point
	}

	//fmt.Println(DIF)
	return result
}

/*
[BERITH]
예상 BC 선출 비율을 계산하는 함수
*/
func (cs *Candidates) getJoinRatio(address common.Address) float64 {
	stake := uint64(0)
	for _, c := range cs.selections {
		if c.address == address {
			stake = c.point
			break
		}
	}

	f := float64(stake) / float64(cs.total)
	r := math.Round(f * float64(100))
	return r
}

type Range struct {
	min   uint64
	max   uint64
	start int
	end   int
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

type VoteResult struct {
	Score *big.Int `json:"score"`
	Rank  int      `json:"rank"`
}

type VoteResults map[common.Address]VoteResult

type sortableList []common.Address

func (s sortableList) Len() int {
	return len(s)
}

func (s sortableList) Swap(a, b int) {
	s[a], s[b] = s[b], s[a]
}

func (s sortableList) Less(a, b int) bool {
	return bytes.Compare(s[a][:], s[b][:]) == -1
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

type JSONCandidate struct {
	Address string `json:"address"`
	Point   uint64 `json:"point"`
	Value   uint64 `json:"value"`
}

type JSONCandidates struct {
	User  []JSONCandidate `json:"user"`
	Total uint64          `json:"total"`
}

func GetCandidates(number uint64, hash common.Hash, stks staking.Stakers, state *state.StateDB) *JSONCandidates {
	list := sortableList(stks.AsList())
	// if len(list) == 0 {
	// 	return result
	// }

	sort.Sort(list)

	cddts := NewCandidates()

	for _, stk := range list {
		point := state.GetPoint(stk).Uint64()
		cddts.Add(Candidate{
			point:   point,
			address: stk,
		})
	}

	jsonCddt := make([]JSONCandidate, 0)
	for _, cddt := range cddts.selections {
		jsonCddt = append(jsonCddt, JSONCandidate{
			Address: cddt.address.Hex(),
			Point:   cddt.point,
			Value:   cddt.val,
		})
	}
	return &JSONCandidates{
		User:  jsonCddt,
		Total: cddts.total,
	}
}

/*
[BERITH]
BC 선출을 하기 위한 함수
선출된 BC map 을 리턴 한다.
*/
func SelectBlockCreator(config *params.ChainConfig, number uint64, hash common.Hash, stks staking.Stakers, state *state.StateDB) VoteResults {
	result := make(VoteResults)

	list := sortableList(stks.AsList())
	if len(list) == 0 {
		return result
	}

	sort.Sort(list)

	cddts := NewCandidates()
	blockNumber := big.NewInt(int64(number))

	for _, stk := range list {
		stakeBalance := state.GetStakeBalance(stk)
		/*
			[Berith]
			Stake Balance 한도 추가에 따라서, Stake Balance를 한도 이상 가지고 있는 대상은 선출 포인트를 재계산
		*/
		var point uint64
		if config.IsBIP4(blockNumber) && stakeBalance.Cmp(config.Bsrr.LimitStakeBalance) == 1  {
			limitStakeBalanceInBer := new(big.Int).Div(config.Bsrr.LimitStakeBalance, big.NewInt(1e+18))
			lastStkBlock := new(big.Int).Set(state.GetStakeUpdated(stk))
			adv := staking.CalcAdvForExceededPoint(blockNumber, lastStkBlock, config.Bsrr.Period, common.BigIntToBigFloat(limitStakeBalanceInBer))

			point = new(big.Int).Add(limitStakeBalanceInBer, adv).Uint64()
		} else {
			point = state.GetPoint(stk).Uint64()
		}

		cddts.Add(Candidate{
			point:   point,
			address: stk,
		})
	}

	if config.IsBIP3(blockNumber) {
		result = cddts.selectBIP3BlockCreator(config, number)
	} else {
		result = cddts.selectBlockCreator(config, number)
	}

	return result
}