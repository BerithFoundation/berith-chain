package staking

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"

	"bitbucket.org/ibizsoftware/berith-chain/consensus"
	"bitbucket.org/ibizsoftware/berith-chain/core/state"

	"bitbucket.org/ibizsoftware/berith-chain/common"
	"bitbucket.org/ibizsoftware/berith-chain/rlp"
)

//var (
//	VoteRatio = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(1))
//)

//StakingMap map implements StakingList
type StakingMap struct {
	storage    map[common.Address]stkInfo
	sortedList []common.Address
}
type stkInfo struct {
	StkAddress     common.Address `json:"address"`
	StkValue       *big.Int       `json:"value"`
	StkBlockNumber *big.Int       `json:"blocknumber"`
}

func (s stkInfo) Address() common.Address { return s.StkAddress }
func (s stkInfo) Value() *big.Int         { return s.StkValue }
func (s stkInfo) BlockNumber() *big.Int   { return s.StkBlockNumber }

func (list *StakingMap) Len() int {
	return len(list.sortedList)
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
		}, nil
	}
	return &stkInfo{
		StkAddress:     address,
		StkValue:       info.Value(),
		StkBlockNumber: info.BlockNumber(),
	}, nil
}

//SetInfo is function to set "staking info"
func (list *StakingMap) SetInfo(info StakingInfo) error {

	if info.Value().Cmp(big.NewInt(0)) < 1 {
		delete(list.storage, info.Address())
	}

	list.storage[info.Address()] = stkInfo{
		StkAddress:     info.Address(),
		StkValue:       info.Value(),
		StkBlockNumber: info.BlockNumber(),
	}
	return nil
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
	for k, v := range list.storage {
		fmt.Println("** [key : ", k.Hex(), " | value : ", v.Value().String(), "| blockNumber : ", v.BlockNumber().String(), "]")
	}
	fmt.Println("==== sortedList ====")
	for _, v := range list.sortedList {
		fmt.Println(v.Hex())
	}
}

//EncodeRLP is function to encode
func (list *StakingMap) EncodeRLP(w io.Writer) error {

	rlpVal, _ := json.Marshal(list.storage)
	//rlpVal[1], _ = json.Marshal(list.sortedList)
	return rlp.Encode(w, rlpVal)
}

func (list *StakingMap) Vote(chain consensus.ChainReader, stateDb *state.StateDB, number uint64, hash common.Hash, epoch uint64, perioid uint64) {
	kv := make(infoForSort, 0)
	for _, v := range list.storage {
		kv = append(kv, v)
	}
	sort.Sort(&kv)

	sortedList := make([]common.Address, 0)

	votes := make([]Vote, 0)
	for _, info := range kv {

		if stateDb == nil {
			break
		}

		reward := stateDb.GetRewardBalance(info.Address())
		v := Vote{info.Address(), info.Value(), info.BlockNumber(), reward}
		votes = append(votes, v)
	}

	if len(votes) > 0 {
		stotal := CalcS(&votes, number, perioid)
		p := CalcP(&votes, stotal, number, perioid)
		r := CalcR(&votes, p)

		//n := common.HexToAddress(header.ParentHash.Hex()).Big().Int64()
		n := common.HexToAddress(hash.Hex()).Big().Int64()
		sig := GetSigners(n, &votes, r, epoch)

		for _, item := range *sig {
			sortedList = append(sortedList, item)
		}
	}

	list.sortedList = sortedList
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
	var btValue []byte
	if err := rlp.DecodeBytes(rlpData, &btValue); err != nil {
		return nil, err
	}

	result := &StakingMap{
		storage:    make(map[common.Address]stkInfo),
		sortedList: make([]common.Address, 0),
	}
	if err := json.Unmarshal(btValue, &result.storage); err != nil {
		return nil, err
	}

	return result, nil

}

//New is function to create new instance
func New() StakingList {
	return &StakingMap{
		storage:    make(map[common.Address]stkInfo),
		sortedList: make([]common.Address, 0),
	}
}
