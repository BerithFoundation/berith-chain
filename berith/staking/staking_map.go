package staking

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"

	"bitbucket.org/ibizsoftware/berith-chain/common"
	"bitbucket.org/ibizsoftware/berith-chain/rlp"
)

//StakingMap map implements StakingList
type StakingMap struct {
	storage    map[common.Address]*big.Int
	sortedList []common.Address
}
type info struct {
	address common.Address
	value   *big.Int
}

func (s *info) Address() common.Address { return s.address }
func (s *info) Value() *big.Int         { return s.value }

func (list *StakingMap) Len() int {
	return len(list.storage)
}

//GetInfoWithIndex is function to get "staking info" that is matched with index from parameter
func (list *StakingMap) GetInfoWithIndex(index int) (StakingInfo, error) {
	if index < 0 || len(list.sortedList) < index {
		return &info{}, errors.New("invalid index")
	}
	address := list.sortedList[index]
	return &info{
		address: address,
		value:   list.storage[address],
	}, nil
}

//GetInfo is function to get "staking info" that is matched with address from parameter
func (list *StakingMap) GetInfo(address common.Address) (StakingInfo, error) {
	value := list.storage[address]

	if value == nil {
		value = big.NewInt(0)
	}
	return &info{
		address: address,
		value:   value,
	}, nil
}

//SetInfo is function to set "staking info"
func (list *StakingMap) SetInfo(address common.Address, x interface{}) error {
	info, ok := x.(*big.Int)

	if ok {
		if info.Cmp(big.NewInt(0)) < 1 {
			delete(list.storage, address)
		}
		list.storage[address] = info
		return nil
	}
	return errors.New("invalid value")
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
		fmt.Println("** [key : ", k.Hex(), " | value : ", v.String(), "]")
	}
}

//EncodeRLP is function to encode
func (list *StakingMap) EncodeRLP(w io.Writer) error {
	rlpVal := make([][]byte, 2)

	rlpVal[0], _ = json.Marshal(list.storage)
	rlpVal[1], _ = json.Marshal(list.sortedList)
	return rlp.Encode(w, rlpVal)
}

func (list *StakingMap) Finalize() {
	list.sort()
}

type infoForSort []info

func (info infoForSort) Len() int           { return len(info) }
func (info infoForSort) Less(i, j int) bool { return info[i].Value().Cmp(info[j].Value()) > 0 }
func (info infoForSort) Swap(i, j int)      { info[i], info[j] = info[j], info[i] }

func (list *StakingMap) sort() {
	kv := make(infoForSort, 0)
	for k, v := range list.storage {
		kv = append(kv, info{address: k, value: v})
	}
	sort.Sort(&kv)

	sortedList := make([]common.Address, 0)
	for _, info := range kv {
		sortedList = append(sortedList, info.Address())
	}

	list.sortedList = sortedList
}

func (list *StakingMap) Copy() *StakingMap {
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
	var btValue [][]byte
	if err := rlp.DecodeBytes(rlpData, &btValue); err != nil {
		return nil, err
	}

	if len(btValue) != 2 {
		return nil, errors.New("failed to get value")
	}

	result := new(StakingMap)
	if err := json.Unmarshal(btValue[0], &result.storage); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(btValue[1], &result.sortedList); err != nil {
		return nil, err
	}

	return result, nil

}

//New is function to create new instance
func New() StakingList {
	return &StakingMap{
		storage:    make(map[common.Address]*big.Int),
		sortedList: make([]common.Address, 0),
	}
}
