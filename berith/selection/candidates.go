package selection

import (
	"crypto/sha256"
	"fmt"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/params"
	"math/big"
	"math/rand"
)

const (
	MAX_MINERS = 10000
)

var (
	maxElectScore = int64(5000000)
	minElectScore = int64(10000)
)

type Candidates struct {
	selections []Candidate
	total      uint64 // Total Selection Point: Staking  + Advantage
	ts         uint64
}

type JSONCandidates struct {
	User  []JSONCandidate `json:"user"`
	Total uint64          `json:"total"`
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
Function to register Staker to elect Block Creator
The function to be called later is the BlockCreator function.
*/
func (cs *Candidates) Add(c Candidate) {
	cs.total += c.point
	c.val = cs.total
	cs.selections = append(cs.selections, c)
}

/*
[Berith]
The block constructor is selected and the result is returned in VoteResults.
*/
func (cs *Candidates) selectBlockCreator(config *params.ChainConfig, number uint64) VoteResults {
	candidateCount := len(cs.selections)
	queue := new(Queue).setQueueAsCandidates(candidateCount)
	result := make(VoteResults)

	currentElectScore := maxElectScore
	electScoreGap := (maxElectScore - minElectScore) / int64(candidateCount)

	// Block number is used as a seed so that all nodes have the same random value
	rand.Seed(cs.GetSeed(config, number))

	err := queue.enqueue(Range{
		min:   0,
		max:   cs.total,
		start: 0,
		end:   candidateCount,
	})
	if err != nil {
		fmt.Println(err)
		return result
	}

	for count := 1; count <= MAX_MINERS && queue.front != queue.rear; count++ {
		r, err := queue.dequeue()
		if err != nil {
			fmt.Println(err)
			return result
		}
		account := r.binarySearch(queue, cs)
		result[account] = VoteResult{
			Score: big.NewInt(currentElectScore + int64(cs.ts)),
			Rank:  count,
		}
		currentElectScore -= electScoreGap
	}
	return result
}

/*
[Berith]
The block constructor is selected and the result is returned in VoteResults.
*/
func (cs *Candidates) selectBIP3BlockCreator(config *params.ChainConfig, number uint64) VoteResults {
	result := make(VoteResults)

	currentElectScore := maxElectScore
	electScoreGap := (maxElectScore - minElectScore) / int64(len(cs.selections))
	rank := 1

	// Block number is used as a seed so that all nodes have the same random value
	rand.Seed(cs.GetSeed(config, number))

	for len(cs.selections) > 0 {
		// The random number below the total elected point is taken and used as the number to select the elected person.
		electedNumber := uint64(rand.Int63n(int64(cs.total)))

		// Search for candidates corresponding to electedNumber by binary search.
		var chosen int
		start := 0
		end := len(cs.selections) - 1
		for {
			mid := (start + end) / 2
			startElectRange := uint64(0)
			if mid > 0 {
				startElectRange = cs.selections[mid-1].val
			}
			endElectRange := cs.selections[mid].val

			if electedNumber >= startElectRange && electedNumber <= endElectRange {
				chosen = mid
				cddt := cs.selections[mid]
				result[cddt.address] = VoteResult{
					Rank:  rank,
					Score: big.NewInt(currentElectScore),
				}
				currentElectScore -= electScoreGap
				rank++
				break
			}

			if electedNumber < startElectRange {
				end = mid - 1
			}
			if electedNumber > endElectRange {
				start = mid + 1
			}
		}

		// Prepare for the selection of next-ranked candidates,
		// except for the data of candidates already elected.
		out := cs.selections[chosen]
		for i := chosen; i+1 < len(cs.selections); i++ {
			newCddt := cs.selections[i+1]
			newCddt.val -= out.point
			cs.selections[i] = newCddt
		}
		cs.selections = cs.selections[:len(cs.selections)-1]
		cs.total -= out.point
	}
	return result
}

/*
[BERITH]
Function to convert block number to hash and force it to int64
Write the result value as Seed.
*/
func (cs Candidates) GetSeed(config *params.ChainConfig, number uint64) int64 {
	// [Berith]
	// Prior to IsBIP2, only 1 byte of the block number is used as a seed
	// After IsBIP2, the entire block number is used as a seed
	bt := []byte{byte(number)}
	if config.IsBIP2(new(big.Int).SetUint64(number)) {
		bt = big.NewInt(0).SetUint64(number).Bytes()
	}

	hash := sha256.New()
	hash.Write(bt)
	md := hash.Sum(nil)
	h := common.BytesToHash(md)
	seed := h.Big().Int64()

	return seed
}