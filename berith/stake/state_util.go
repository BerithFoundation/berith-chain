package stake

import (
	"bytes"
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

func NewStakingInfo(address *common.Address, value *big.Int) StakingInfo {
	return &info {
		address : *address,
		value : value,
	}
}

func (s *info) Address() common.Address { return s.address }
func (s *info) Value() *big.Int         { return s.value }

func (list *StakingMap) getAddressWithIndex(index int) (common.Address, error) {
	if index < 0 || len(list.sortedList) < index {
		return common.Address{}, errors.New("invalid index")
	}
	return list.sortedList[index], nil
}

func (list *StakingMap) Len() int {
	return len(list.storage)
}

func (list *StakingMap) getIndexWithAddress(address common.Address) (int, error) {
	for i, value := range list.sortedList {
		if bytes.Compare(value.Bytes(), address.Bytes()) == 0 {
			return i, nil
		}
	}
	return -1, errors.New("no index matched")
}

func (list *StakingMap) getMinerWithAddressAndNum(address common.Address, num int) (common.Address, error) {
	index, err := list.getIndexWithAddress(address)
	if err != nil {
		return common.Address{}, err
	}

	result, getErr := list.getAddressWithIndex(index + num)
	if getErr != nil {
		return common.Address{}, getErr
	}

	if bytes.Compare(result.Bytes(), common.Address{}.Bytes()) == 0 {
		return common.Address{}, errors.New("no value matched")
	}

	return result, nil
}

func (list *StakingMap) NextMiner(address common.Address) (common.Address, error) {
	return list.getMinerWithAddressAndNum(address, 1)
}

func (list *StakingMap) PrevMiner(address common.Address) (common.Address, error) {
	return list.getMinerWithAddressAndNum(address, -1)
}

func (list *StakingMap) GetMiner(index int) (common.Address, error) {
	result, err := list.getAddressWithIndex(index)

	if err != nil {
		return common.Address{}, err
	}

	if bytes.Compare(result.Bytes(), common.Address{}.Bytes()) == 0 {
		return common.Address{}, errors.New("no value matched")
	}

	return result, nil
}

//Get getter of StakingMap
func (list *StakingMap) Get(address common.Address) (StakingInfo, error) {
	value := list.storage[address]

	if value == nil {
		value = big.NewInt(0)
		return nil, nil
	}
	
	return &info{
		address: address,
		value:   value,
	}, nil
}

//Set setter of StakingMap
func (list *StakingMap) Set(address common.Address, x interface{}) error {
	info, ok := x.(*big.Int)

	if ok {
		list.storage[address] = info
		return nil
	}
	return errors.New("invalid value")
}

// UnSet address from the staking list
func (list *StakingMap) Delete(address common.Address) error {
	if _, ok := list.storage[address]; ok {
		delete(list.storage, address)
	}
	return nil
}

// Print is
func (list *StakingMap) Print() {
	fmt.Println("==== Staking List ====")
	for k, v := range list.storage {
		fmt.Println("** [key : ", k.Hex(), " | value : ", v.String(), "]")
	}
}

func (list *StakingMap) EncodeRLP(w io.Writer) error {
	rlpVal := make([][]byte, 2)

	rlpVal[0], _ = json.Marshal(list.storage)
	rlpVal[1], _ = json.Marshal(list.sortedList)
	return rlp.Encode(w, rlpVal)
}

//NewStakingList get staking list to trie
func NewStakingMap(db DataBase, blockNumber *big.Int, hash common.Hash) (*StakingMap, error) {
	if blockNumber.Cmp(big.NewInt(0)) <= 0 {
		return &StakingMap{
			storage:    make(map[common.Address]*big.Int, 0),
			sortedList: make([]common.Address, 0),
		}, nil
	}

	rlpData, err1 := db.GetValue(hash.Hex() + ":" + blockNumber.String())
	if err1 != nil {
		return nil, err1
	}

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

//Commit commit staking list to database
func (list *StakingMap) Commit(db DataBase, blockNumber *big.Int, hash common.Hash) error {
	list.sort()

	rlpValue, err := rlp.EncodeToBytes(list)

	if err != nil {
		return err
	}
	db.PushValue(hash.Hex()+":"+blockNumber.String(), rlpValue)

	return nil
}
