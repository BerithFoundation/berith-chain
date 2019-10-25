// Copyright 2019 The berith Authors
// This file is part of berith.
//
// berith is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// berith is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with berith. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/BerithFoundation/berith-chain/berithclient"
	cli "gopkg.in/urfave/cli.v1"
	"math/big"
	"math/rand"
	"strings"
)

var (
	// flags
	StartBlockFlag = cli.Int64Flag{
		Name:  "startblock",
		Usage: "first block number to check",
		Value: int64(1),
	}
	LastBlockFlag = cli.Int64Flag{
		Name:  "lastblock",
		Usage: "last block number to check",
		Value: int64(0),
	}

	// commands
	TpsCommand = cli.Command{
		Action: testTPS,
		Name:   "tps",
		Usage:  "test tps",
		Flags: []cli.Flag{
			NodesFlag,
			StartBlockFlag,
			LastBlockFlag,
		},
		Description: `this is tps command.`,
	}
)

func testTPS(ctx *cli.Context) error {
	cctx := context.Background()
	// parse nodes
	nodes := strings.Split(ctx.String(NodesFlag.Name), ",")
	if len(nodes) == 1 && nodes[0] == "" {
		return errors.New("must have at least one node rpc url")
	}
	node := nodes[rand.Intn(100)%len(nodes)]
	client, err := berithclient.Dial(node)
	if err != nil {
		return err
	}
	defer client.Close()

	// parse block range
	startNumber := ctx.Int64(StartBlockFlag.Name)
	if startNumber < 1 {
		return errors.New("start block number must larger than 0")
	}
	lastNumber := ctx.Int64(LastBlockFlag.Name)
	lastBlock, err := client.BlockByNumber(cctx, nil)
	if err != nil {
		return err
	}
	lastBlockNumber := lastBlock.Number().Int64()
	if lastNumber == 0 || lastNumber > lastBlockNumber {
		lastNumber = lastBlockNumber
	}
	if lastNumber < startNumber {
		return errors.New("last block number must be larger than start block number")
	}

	fmt.Printf(">>>>>> start to test tps with range (%d,%d] <<<<<<\n", startNumber, lastNumber)
	var startBlockTime, lastBlockTime, lastTime, total int64
	for blockNumber := startNumber; blockNumber <= lastNumber; blockNumber++ {
		block, err := client.BlockByNumber(cctx, big.NewInt(blockNumber))
		if err != nil {
			return err
		}
		current := block.Time().Int64()
		if blockNumber == startNumber {
			startBlockTime = current
			lastTime = current
			continue
		}
		if blockNumber == lastNumber {
			lastBlockTime = current
		}

		txLen := int64(len(block.Transactions()))
		total += txLen
		tps := float64(txLen) / float64(current-lastTime)
		fmt.Printf("block #%d > #tx : %d / #tps : %.4f\n", blockNumber, txLen, tps)
		lastTime = current
	}

	totalTps := float64(total) / float64(lastBlockTime-startBlockTime)
	fmt.Printf(">> block (%d,%d] #tps : %.4f / #tx : %d\n", startNumber, lastNumber, totalTps, total)
	return nil
}
