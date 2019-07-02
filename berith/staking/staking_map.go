package staking

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/rlp"
)

//var (
//	VoteRatio = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(1))
//)

var (
	roi		map[common.Address]float64
)

//StakingMap map implements StakingList
type StakingMap struct {
	storage    map[common.Address]stkInfo
	sortedList []common.Address
	miners     map[common.Address]bool
	table      map[common.Address]*big.Int
	target     common.Hash
}
type stkInfo struct {
	StkAddress     common.Address `json:"address"`
	StkValue       *big.Int       `json:"value"`
	StkBlockNumber *big.Int       `json:"blocknumber"`
	StkReward      *big.Int       `json:"reward"`
}

func (s stkInfo) Address() common.Address { return s.StkAddress }
func (s stkInfo) Value() *big.Int         { return s.StkValue }
func (s stkInfo) BlockNumber() *big.Int   { return s.StkBlockNumber }
func (s stkInfo) Reward() *big.Int        { return s.StkReward }

func (list *StakingMap) SetTarget(target common.Hash) {
	list.target = target
}

func (list *StakingMap) GetTarget() common.Hash {
	return list.target
}

func (list *StakingMap) Len() int {
	return len(list.sortedList)
}

func (list *StakingMap) SetMiner(address common.Address) {
	list.miners[address] = true
}

func (list *StakingMap) InitMiner() {
	list.miners = make(map[common.Address]bool)
}

func (list *StakingMap) GetMiners() map[common.Address]bool {
	return list.miners
}

func (list *StakingMap) GetDifficulty(addr common.Address, blockNumber, period uint64) (*big.Int, bool) {
	flag := false
	if len(list.table) <= 0 {
		flag = true
		list.selectSigner(blockNumber, period)
	}
	if len(list.table) <= 0 {
		return big.NewInt(1234), false
	}

	result, ok := list.table[addr]
	if !ok {
		result = big.NewInt(0)
	}
	return result, flag
}

//GetInfoWithIndex is function to get "staking info" that is matched with index from parameter
func (list *StakingMap) GetInfoWithIndex(index int) (StakingInfo, error) {
	if index < 0 || len(list.sortedList) < index {
		return stkInfo{}, errors.New("invalid index")
	}
	address := list.sortedList[index]
	return list.storage[address], nil
}

//GetInfo is function to get "staking info" that is matched with address from parameter
func (list *StakingMap) GetInfo(address common.Address) (StakingInfo, error) {
	info, ok := list.storage[address]

	if !ok {
		return &stkInfo{
			StkAddress:     address,
			StkValue:       big.NewInt(0),
			StkBlockNumber: big.NewInt(0),
			StkReward:      big.NewInt(0),
		}, nil
	}
	return &stkInfo{
		StkAddress:     address,
		StkValue:       info.Value(),
		StkBlockNumber: info.BlockNumber(),
		StkReward:      info.Reward(),
	}, nil
}

//SetInfo is function to set "staking info"
func (list *StakingMap) SetInfo(info StakingInfo) error {

	if info.Value().Cmp(big.NewInt(0)) < 1 && info.Reward().Cmp(big.NewInt(0)) < 1 {
		delete(list.storage, info.Address())
	}

	list.storage[info.Address()] = stkInfo{
		StkAddress:     info.Address(),
		StkValue:       info.Value(),
		StkBlockNumber: info.BlockNumber(),
		StkReward:      info.Reward(),
	}
	return nil
}

func (list *StakingMap) ToArray() []common.Address {
	return list.sortedList
}

//Delete is function to delete address from the staking list
func (list *StakingMap) Delete(address common.Address) error {
	if _, ok := list.storage[address]; ok {
		delete(list.storage, address)
	}
	return nil
}

// Print is function to print stakingList info
func (list *StakingMap) Print() {
	fmt.Println("==== Staking List ====")
	fmt.Println("TARGET :", list.target.Hex())
	for k, v := range list.storage {
		fmt.Println("** [key : ", k.Hex(), " | value : ", v.Value().String(), "| blockNumber : ", v.BlockNumber().String(), "| reward : ", new(big.Int).Div(v.Reward(), big.NewInt(1000000000000000000)), "]")
	}
	fmt.Println("==== sortedList ====")
	for _, v := range list.sortedList {
		fmt.Println(v.Hex())
	}
	// fmt.Println("====== MINERS ======")
	// for k, v := range list.miners {
	// 	fmt.Println("[", k.Hex(), ",", v, "]")
	// }
}

//EncodeRLP is function to encode
func (list *StakingMap) EncodeRLP(w io.Writer) error {

	var byteArr [5][]byte

	byteArr[0], _ = json.Marshal(list.storage)
	byteArr[1], _ = json.Marshal(list.miners)
	byteArr[2] = list.target[:]
	byteArr[3], _ = json.Marshal(list.table)
	byteArr[4], _ = json.Marshal(list.sortedList)
	//rlpVal[1], _ = json.Marshal(list.sortedList)
	return rlp.Encode(w, byteArr)
}

func (list *StakingMap) Sort() {

	if len(list.sortedList) > 0 {
		return
	}

	kv := make(infoForSort, 0)
	for _, v := range list.storage {
		if v.Value().Cmp(big.NewInt(0)) > 0 {
			kv = append(kv, v)
		}
	}
	sort.Sort(&kv)

	sortedList := make([]common.Address, 0)

	for _, info := range kv {

		sortedList = append(sortedList, info.Address())
	}

	list.sortedList = sortedList
}

func (list *StakingMap) ClearTable() {
	list.sortedList = make([]common.Address, 0)
	list.table = make(map[common.Address]*big.Int)
}

func (list *StakingMap) selectSigner(blockNumber, period uint64) {

	if len(list.sortedList) <= 0 {
		list.Sort()
	}

	if len(list.sortedList) <= 0 {
		return
	}

	//cs := &Candidates{
	//	number:     blockNumber,
	//	period:     period,
	//	selections: make([]Candidate, 0),
	//}

	cs := NewCandidates(blockNumber, period)

	for _, addr := range list.sortedList {
		info := list.storage[addr]
		reward := info.StkReward
		if reward == nil {
			reward = big.NewInt(0)
		}
		value, _ := new(big.Int).SetString(info.Value().String(), 10)
		blockNumber, _ := new(big.Int).SetString(info.BlockNumber().String(), 10)
		cs.Add(Candidate{info.Address(), value, blockNumber, reward})

	}

	list.table = *cs.GetBlockCreator(blockNumber)
	roi = cs.GetRoi()


	for key, value := range list.table {
		fmt.Println("ADDRESS :: "+key.String(), "DIFF :: "+value.String())
	}

}

func (list *StakingMap) GetRoi(address common.Address) float64 {
	return roi[address]
}

type infoForSort []stkInfo

func (info infoForSort) Len() int { return len(info) }
func (info infoForSort) Less(i, j int) bool {
	if info[i].Value().Cmp(info[j].Value()) == 0 {
		if info[i].BlockNumber().Cmp(info[j].BlockNumber()) == 0 {
			return bytes.Compare(info[i].Address().Bytes(), info[j].Address().Bytes()) > 0
		}
		return info[i].BlockNumber().Cmp(info[j].BlockNumber()) < 0
	}
	return info[i].Value().Cmp(info[j].Value()) > 0
}
func (info infoForSort) Swap(i, j int) { info[i], info[j] = info[j], info[i] }

func (list *StakingMap) Copy() StakingList {
	return &StakingMap{
		storage:    list.storage,
		sortedList: list.sortedList,
		miners:     list.miners,
		target:     list.target,
		table:      list.table,
	}
}

func (list *StakingMap) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(list)
}

func (list *StakingMap) Decode(rlpData []byte) (StakingList, error) {
	return Decode(rlpData)
}

func Encode(stakingList StakingList) ([]byte, error) {
	return rlp.EncodeToBytes(stakingList)
}

func Decode(rlpData []byte) (StakingList, error) {
	var byteArr [5][]byte
	if err := rlp.DecodeBytes(rlpData, &byteArr); err != nil {
		return nil, err
	}

	result := &StakingMap{
		storage:    make(map[common.Address]stkInfo),
		sortedList: make([]common.Address, 0),
		miners:     make(map[common.Address]bool),
		table:      make(map[common.Address]*big.Int),
		target:     common.Hash{},
	}
	if err := json.Unmarshal(byteArr[0], &result.storage); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteArr[1], &result.miners); err != nil {
		return nil, err
	}

	result.target = common.BytesToHash(byteArr[2])

	if err := json.Unmarshal(byteArr[3], &result.table); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteArr[4], &result.sortedList); err != nil {
		return nil, err
	}

	return result, nil

}

//New is function to create new instance
func New() StakingList {
	return &StakingMap{
		storage:    make(map[common.Address]stkInfo),
		sortedList: make([]common.Address, 0),
		miners:     make(map[common.Address]bool),
		table:      make(map[common.Address]*big.Int),
	}
}
