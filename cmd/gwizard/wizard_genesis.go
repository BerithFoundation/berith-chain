// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"bitbucket.org/ibizsoftware/berith-chain/common"
	"bitbucket.org/ibizsoftware/berith-chain/core"
	"bitbucket.org/ibizsoftware/berith-chain/log"
	"bitbucket.org/ibizsoftware/berith-chain/params"
)

// makeGenesis creates a new genesis struct based on some user input.
func (w *wizard) makeGenesis() {
	// Construct a default genesis block
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   4700000,
		Difficulty: big.NewInt(524288),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			HomesteadBlock:      big.NewInt(1),
			EIP150Block:         big.NewInt(2),
			EIP155Block:         big.NewInt(3),
			EIP158Block:         big.NewInt(3),
			ByzantiumBlock:      big.NewInt(4),
			ConstantinopleBlock: big.NewInt(5),
		},
	}

	// In the case of clique, configure the consensus parameters
	genesis.Difficulty = big.NewInt(1)
	genesis.Config.Bsrr = &params.BSRRConfig{
		Period: 30,
		Epoch:  300,
		Rewards:  big.NewInt(500),
		StakeMinimum: new(big.Int).Mul(big.NewInt(100000), big.NewInt(1e+18)),
		SlashRound: uint64(1),
	}

	fmt.Println()
	fmt.Println("What is network name?")
	w.network = w.readDefaultString("Genesis")

	fmt.Println()
	fmt.Println("How many seconds should blocks take? (default = 15)")
	genesis.Config.Bsrr.Period = uint64(w.readDefaultInt(15))

	// We also need the initial list of signers
	fmt.Println()
	fmt.Println("Which accounts are allowed to seal (First Block Creator)?")

	var signers []common.Address
	if address := w.readAddress(); address != nil {
		signers = append(signers, *address)
	}


	// Sort the signers and embed into the extra-data section
	for i := 0; i < len(signers); i++ {
		for j := i + 1; j < len(signers); j++ {
			if bytes.Compare(signers[i][:], signers[j][:]) > 0 {
				signers[i], signers[j] = signers[j], signers[i]
			}
		}
	}

	genesis.ExtraData = make([]byte, 32+len(signers)*common.AddressLength+65)
	for i, signer := range signers {
		copy(genesis.ExtraData[32+i*common.AddressLength:], signer[:])
	}


	//Set First BC Balance
	if address := &signers[0]; address != nil {
		genesis.Alloc[*address] = core.GenesisAccount{
			Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
		}
	}

	// Query the user for some custom extras
	fmt.Println()
	fmt.Println("Specify your chain/network ID if you want an explicit one (default = random)")
	genesis.Config.ChainID = new(big.Int).SetUint64(uint64(w.readDefaultInt(rand.Intn(65536))))

	w.conf.Genesis = genesis
	w.conf.flush()


	folder := w.readDefaultString(".")
	if err := os.MkdirAll(folder, 0755); err != nil {
		log.Error("Failed to create spec folder", "folder", folder, "err", err)
		return
	}
	out, _ := json.MarshalIndent(w.conf.Genesis, "", "  ")

	// Export the native genesis spec used by puppeth and Geth
	json := filepath.Join(folder, fmt.Sprintf("%s.json", w.network))
	if err := ioutil.WriteFile((json), out, 0644); err != nil {
		log.Error("Failed to save genesis file", "err", err)
		return
	}
	log.Info("Saved native genesis chain spec", "path", json)

	// Export the genesis spec used by Harmony (formerly EthereumJ
	saveGenesis(folder, w.network, "harmony", w.conf.Genesis)
}



// saveGenesis JSON encodes an arbitrary genesis spec into a pre-defined file.
func saveGenesis(folder, network, client string, spec interface{}) {
	path := filepath.Join(folder, fmt.Sprintf("%s-%s.json", network, client))

	out, _ := json.Marshal(spec)
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		log.Error("Failed to save genesis file", "client", client, "err", err)
		return
	}
	log.Info("Saved genesis chain spec", "client", client, "path", path)
}
