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
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/accounts/keystore"
	"github.com/BerithFoundation/berith-chain/berithclient"
	"github.com/BerithFoundation/berith-chain/cmd/utils"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/log"
	cli "gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	// flags
	KeystoreFlag = cli.StringFlag{
		Name:  "keystore",
		Usage: "Directory of keystore file",
	}
	AddressesFlag = cli.StringFlag{
		Name:  "addresses",
		Usage: "From addresses to send transactions.",
	}
	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Password file path",
	}
	DurationFlag = cli.StringFlag{
		Name:  "duration",
		Usage: "How long the test will be executed",
	}
	TxCountFlag = cli.Uint64Flag{
		Name:  "txcount",
		Usage: "How many test runs will be executed",
	}
	TxIntervalFlag = cli.Uint64Flag{
		Name:  "txinterval",
		Usage: "Interval between transactions [ms]",
	}
	InitDelay = cli.Uint64Flag{
		Name:  "initdelay",
		Usage: "Sleep before testing",
	}
	OutputPath = cli.StringFlag{
		Name:  "outputpath",
		Usage: "Dir path of output",
	}

	// commands
	ExecuteCommand = cli.Command{
		Action:   cli.ShowSubcommandHelp,
		Name:     "execute",
		Usage:    "request transactions to given nodes",
		Category: "EXECUTE COMMANDS",
		Subcommands: []cli.Command{
			{
				Name:   "transfer",
				Usage:  "sends only transfer transactions",
				Action: transfer,
				Flags: []cli.Flag{
					ChainIdFlag,
					NodesFlag,
					ConfigFileFlag,
					KeystoreFlag,
					AddressesFlag,
					PasswordFlag,
					DurationFlag,
					TxCountFlag,
					TxIntervalFlag,
					InitDelay,
					OutputPath,
				},
			},
		},
	}
)

type doWork func(taskID string, count uint64)

// transfer send transactions given cli context
func transfer(ctx *cli.Context) error {
	// setup
	cfg, err := parseConfig(ctx)
	if err != nil {
		return err
	}

	testCtx := parseContext(cfg)
	if testCtx.TxCount < 1 && testCtx.Duration <= 0 {
		utils.Fatalf("invalid tx count %d with duration %s", testCtx.TxCount, testCtx.Duration.String())
	}

	setupContext(testCtx)
	defer tearDownContext(testCtx)

	// do request after initial delay
	if testCtx.InitDelay > 0 {
		time.Sleep(time.Duration(testCtx.InitDelay))
	}
	log.Info(">> Start to send transfer transactions <<")
	startTime := time.Now()
	var wait sync.WaitGroup
	for i, addrCtx := range testCtx.AddressContexts {
		wait.Add(1)
		taskID := "task-" + strconv.Itoa(i)
		go func(id string, a *addressContext, w *sync.WaitGroup) {
			sendTransferTransactions(id, testCtx, a, w)
		}(taskID, addrCtx, &wait)
	}
	wait.Wait()
	endTime := time.Now()

	// TODO : write output to console and file by using text/template
	layout := "15:04:05.999"
	var out bytes.Buffer
	out.WriteString("// ------------------------------------\n")
	out.WriteString(fmt.Sprintf(">>>>>> Complete to send transactions [%s] ~ [%s] <<<<<<\n", startTime.Format(layout), endTime.Format(layout)))

	out.WriteString("## args\n")
	out.WriteString(fmt.Sprintf("> repeat : [%d], duration : [%s], interval : [%dms]\n", testCtx.TxCount, testCtx.Duration, testCtx.TxInterval))

	out.WriteString("\n## nodes\n")
	for _, nodeCtx := range testCtx.NodeContexts {
		out.WriteString(fmt.Sprintf("Node : %s --> %s\n", nodeCtx.url, nodeCtx.summary.String()))
	}

	out.WriteString("\n## addrs\n")
	for _, addrCtx := range testCtx.AddressContexts {
		out.WriteString(fmt.Sprintf("Addr : %s --> %s, first tx hash : %s\n", addrCtx.account.Address.Hex(), addrCtx.summary.String(), addrCtx.summary.firstTxHash))
	}
	out.WriteString("--------------------------------------- //\n")

	// flush output to console & file
	fmt.Print(out.String())
	if testCtx.OutputPath != "" {
		path := testCtx.OutputPath
		fi, err := os.Stat(path)
		if err == nil {
			// if exist path, must be directory.
			if !fi.IsDir() {
				log.Warn("already exist output path.", "path", path)
				return nil
			}
			path = filepath.Join(path, "berithench-"+startTime.Format(layout)+".txt")
		} else if !os.IsNotExist(err) {
			log.Warn("failed to write output", "err", err)
			return nil
		}
		fmt.Println(">> try to write output to ", path)
		err = ioutil.WriteFile(path, out.Bytes(), 0644)
		if err != nil {
			log.Warn("failed to write a output file", "err", err)
			return nil
		}
	}
	return nil
}

// sendTransferTransactions request transactions given test context with an account.
func sendTransferTransactions(taskID string, ctx *berithenchContext, addrCtx *addressContext, wait *sync.WaitGroup) {
	var toAddrs []string
	// if only one account, then add addresses from 1 to 9
	if len(ctx.AddressContexts) == 1 {
		toAddrs = []string{
			"0000000000000000000000000000000000000001",
			"0000000000000000000000000000000000000002",
			"0000000000000000000000000000000000000003",
			"0000000000000000000000000000000000000004",
			"0000000000000000000000000000000000000005",
			"0000000000000000000000000000000000000006",
			"0000000000000000000000000000000000000007",
			"0000000000000000000000000000000000000008",
			"0000000000000000000000000000000000000009",
		}
	}
	// extract to addresses's hex
	for _, atx := range ctx.AddressContexts {
		if atx == addrCtx {
			continue
		}
		toAddrs = append(toAddrs, atx.account.Address.Hex()[2:])
	}
	fmt.Printf("final to addrs :%v\n", toAddrs)

	// setup task
	cctx := context.Background() // TODO : check need timeout
	work := func(taskID string, taskCount uint64) {
		// taskCount is started with 1
		nonce := addrCtx.startNonce + taskCount - 1
		randIdx := rand.Intn(100) % len(toAddrs)
		to := common.HexToAddress(toAddrs[randIdx])
		gas := uint64(21000)
		gasPrice := common.Big1
		value := *common.Big1
		tx := types.NewTransaction(
			nonce,
			to,
			&value,
			gas,
			gasPrice,
			nil,
			types.Main,
			types.Main,
		)

		tx, err := ctx.Keystore.SignTx(addrCtx.account, tx, ctx.ChainID)
		if err != nil {
			log.Warn("cannot sign a transaction", "error", err)
			return
		}
		randIdx = rand.Intn(100) % len(ctx.NodeContexts)
		nodeCtx := ctx.NodeContexts[randIdx]
		err = nodeCtx.client.SendTransaction(cctx, tx)
		success := err == nil
		addrCtx.summary.addResult(success)
		nodeCtx.summary.addResult(success)
		if err != nil {
			// TODO : add fail queue
			fmt.Println(err)
		}
		if err == nil && taskCount == 1 {
			addrCtx.summary.firstTxHash = tx.Hash().Hex()
		}
	}
	var repeat uint64
	if addrCtx.lastNonce == 0 {
		repeat = 0
	} else {
		repeat = addrCtx.lastNonce - addrCtx.startNonce + 1
	}
	interval := time.Duration(ctx.TxInterval)
	delay := ctx.Duration

	err := schedule(taskID, repeat, interval, delay, wait, work)
	if err != nil {
		log.Error("failed to setup task.", "addr", addrCtx.account.Address.Hex())
	}
}

// setupContext setup test context such as calc nonce.
func setupContext(ctx *berithenchContext) {
	// setup start,end accounts nonce
	addrLen := len(ctx.AddressContexts)
	for i, addrCtx := range ctx.AddressContexts {
		node := ctx.NodeContexts[i%len(ctx.NodeContexts)]
		startNonce, err := node.client.PendingNonceAt(ensureContext(nil), addrCtx.account.Address)
		if err != nil {
			utils.Fatalf("failed to get start nonce. node : %s. error :%v", node.url, err)
		}
		addrCtx.startNonce = startNonce
		if ctx.TxCount < 1 {
			// ta.lastNonce = ^uint64(0)
			addrCtx.lastNonce = 0
		} else {
			txCount := calcTxCount(i, addrLen, int(ctx.TxCount))
			addrCtx.lastNonce = startNonce + uint64(txCount) - 1
		}
	}
}

// tearDownContext clean up test context if needed
func tearDownContext(ctx *berithenchContext) {
	// clear node context
	for _, nodeCtx := range ctx.NodeContexts {
		nodeCtx.close()
	}
}

// calcTxCount calculate counts of sending transactions depends on length of addresses
func calcTxCount(idx, addrLen, txCount int) int {
	if txCount < 1 {
		return 0
	}
	count := txCount / addrLen
	if idx < txCount%addrLen-1 {
		count++
	}
	return count
}

// parseContext parse bbench config to test context
func parseContext(config *berithenchConfig) *berithenchContext {
	ctx := berithenchContext{
		ChainID:    big.NewInt(config.ChainID),
		Duration:   parseDuration(config),
		TxCount:    config.TxCount,
		TxInterval: config.TxInterval,
		InitDelay:  config.InitDelay,
		OutputPath: config.OutputPath,
	}

	// parse nodes rpc url -> testNode
	ctx.NodeContexts = make([]*nodeContext, len(config.Nodes))
	for i, node := range config.Nodes {
		ctx.NodeContexts[i] = NewNodeContext(node)
		client, err := berithclient.Dial(node)
		if err != nil {
			utils.Fatalf("failed to new berith client at node %s\n", node)
		}
		ctx.NodeContexts[i].client = client
	}

	// parse addresses, password -> test addresses and unlock
	targetAddrs := make(map[string]string)
	addrs := config.Addresses
	passwords := makePasswordList(config)
	if passwords == nil {
		utils.Fatalf("must exist password to send transactions.")
	}
	if len(passwords) == 1 {
		for i := 0; i < len(addrs)-1; i++ {
			passwords = append(passwords, passwords[0])
		}
	}
	if len(passwords) != len(addrs) {
		utils.Fatalf("different count between addresses and passwords. addrs : #%d, passwords : #%d", len(addrs), len(passwords))
	}
	for i, addr := range addrs {
		targetAddrs[addr] = passwords[i]
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ctx.Keystore = keystore.NewKeyStore(config.Keystore, scryptN, scryptP)
	// unlock and make test address
	for addr, pass := range targetAddrs {
		account, err := ctx.Keystore.Find(accounts.Account{Address: common.HexToAddress(addr)})
		if err != nil {
			utils.Fatalf("cannot find address %s in keystore.", addr)
		}
		err = ctx.Keystore.Unlock(account, pass)
		if err != nil {
			utils.Fatalf("cannot unlock address %s in keystore.", addr)
		}
		a := NewAddressContext(account)
		ctx.AddressContexts = append(ctx.AddressContexts, a)
	}
	return &ctx
}

// ensureContext is a helper method to ensure a context is not nil, even if the
// user specified it as such.
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		// return context.TODO()
		return context.Background()
	}
	return ctx
}

// schedule execute doWork func repeatedly.
// if repeat is larger than 0, then execute work function of repeat times,
// otherwise execute with delay. task count is always started with 1.
func schedule(taskID string, repeat uint64, interval, delay time.Duration, w *sync.WaitGroup, work doWork) error {
	if repeat < 1 && delay < 1 {
		return errors.New("repeat or delay must be larger than 0")
	}

	// interval must larger than 0 to use ticker
	if interval < 1 {
		interval = 1
	}
	if repeat > 1 {
		delay = math.MaxInt64
	}

	ticker := time.NewTicker(interval)
	cancel := make(chan bool)
	timeout := time.NewTimer(delay)

	count := uint64(1)
	defer w.Done()
	for {
		if repeat > 0 && count > repeat {
			return nil
		}
		select {
		case <-timeout.C:
			return nil
		case <-cancel:
			return nil
		case <-ticker.C:
			work(taskID, count)
			count++
		}
	}
}
