package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/berith/stakingdb"
	"github.com/ethereum/go-ethereum/berith/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/node"
	"math/big"
	"path/filepath"
)

func main()  {
	filepath.Join(node.DefaultDataDir(), "geth", "chaindata")

	sk := &stakingdb.StakingDB{}

	err := sk.CreateDB("stakingDB")
	defer sk.Close()

	if err != nil {
		fmt.Println(err.Error())
		return
	}


	//GET
	s, err := sk.GetValue(berith.StakingListKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	addr := []byte("ADDRESS")
	acc := common.BytesToAddress(addr)

	if s != nil {
		fmt.Println("GET VALUE :: ", s[acc])
	}



	//PUT
	var value = make(map[common.Address]*big.Int)
	value[acc] = big.NewInt(1000)

	err = sk.PushValue(berith.StakingListKey, value)
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	fmt.Println("SUCCESS")

}

