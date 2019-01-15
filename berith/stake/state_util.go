package stake

import (

	"bitbucket.org/ibizsoftware/berith-chain/common/stakesort"
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
	storage *map[common.Address]*big.Int
	stakinglist *stakesort.Stakelist

}
type storage struct {
	address common.Address
	value   *big.Int
}

func (s storage) Address() common.Address { return s.address }
func (s storage) Value() *big.Int         { return s.value }
func (list StakingMap) GetRRList() *stakesort.Stakelist{
	return list.stakinglist
}
//Get getter of StakingMap
func (list StakingMap) Get(address common.Address) (StakingInfo, error) {
	value := list.storage
	x:=value[address]

	if value == nil {
		value = big.NewInt(0)
	}
	return storage{
		address: address,
		value:   value,
	}, nil
}

//Set setter of StakingMap
func (list StakingMap) Set(address common.Address, x interface{}) error {
	info, ok := x.(*big.Int)
	if ok {
		list.storage[address] = info
		return nil
	} else {
		return errors.New("invalid value")
	}
}

func (list StakingMap) Copy() StakingList {
	fmt.Println("22222-111111")
	return StakingMap{
		storage: list.storage,
	}
}

func (list StakingMap) Print() {
	fmt.Println(list)
}

func (list StakingMap) EncodeRLP(w io.Writer) error {
	result, _ := json.Marshal(list.storage)
	return rlp.Encode(w, result)
}

//GetStakingList get staking list to trie
func GetStakingMap(db DataBase, blockNumber *big.Int, hash common.Hash) (StakingList, error) {
	rlpData, err1 := db.GetValue(hash.Hex() + ":" + blockNumber.String())
	if err1 != nil {
		return nil, nil
	}

	var btValue []byte
	if err := rlp.DecodeBytes(rlpData, &btValue); err != nil {
		return nil, err
	}

	var result map[common.Address]*big.Int
	if err := json.Unmarshal(btValue, &result); err != nil {
		return nil, err
	}
	var stakelist stakesort.Stakelist
	for addr,value := range result{
		stakelist = append(stakelist,&stakesort.Stake{addr,value})
	}
	sort.Sort(stakelist)

	return &StakingMap{
		storage: &result,
		stakinglist: &stakelist,
	}, nil
}

//Commit commit staking list to database
func Commit(list StakingList, db DataBase, blockNumber *big.Int, hash common.Hash) error {
	rlpValue, err := rlp.EncodeToBytes(list)

	if err != nil {
		return err
	}

	db.PushValue(hash.Hex()+":"+blockNumber.String(), rlpValue)

	return nil

}

func AppendTransaction(list StakingList, tx Transaction) error {
	if !tx.Staking() {
		return errors.New("not staking Transaction")
	}
	info, getErr := list.Get(tx.From())
	if getErr != nil {
		return getErr
	}

	if err := list.Set(tx.From(), big.NewInt(0).Add(info.Value(), tx.Value())); err != nil {
		return err
	}

	return nil
}
