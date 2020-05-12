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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"berith-chain/accounts"
	"berith-chain/accounts/keystore"
	"berith-chain/berithclient"
	"berith-chain/cmd/utils"
	"berith-chain/common"
	"berith-chain/common/hexutil"
	"berith-chain/core/types"
	"berith-chain/rpc"
	"github.com/gookit/color"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	transferResultTemplate = `// ------------------------------------------------------------------------
>>>>>> Complete to send transactions [{{.StartTime}}] ~ [{{.EndTime}}] <<<<<<
## args
> repeat : [{{.Repeat}}], duration : [{{.Duration}}], interval : [{{.Interval}}]
## nodes
{{range $index, $node := .Nodes}}[#{{$index}}] Node : {{$node.Url}} --> Request : {{$node.Request}} / Success : {{$node.Success}} / Fail : {{$node.Fail}}
{{end}}
## addresses
{{range $index, $address := .Addresses}}[#{{$index}}] Addr({{$address.Hex}})
 - Request : {{$address.Request}} / Success : {{$address.Success}} / Fail : {{$address.Fail}}
 - first tx hash : {{$address.FirstTxHash}} (in block {{$address.FirstTxBlockNumber}})
{{end}}
--------------------------------------------------------------------------- //`
)

var (
	// flags

	AddressesFlag = cli.StringFlag{
		Name:  "addresses",
		Usage: "From addresses to send transactions.",
	}

	DurationFlag = cli.StringFlag{
		Name:  "duration",
		Usage: "How long the test will be executed",
	}

	InitDelay = cli.Uint64Flag{
		Name:  "initdelay",
		Usage: "Sleep before testing",
	}
	EnableCpuProfile = cli.BoolFlag{
		Name:  "cpuprofile",
		Usage: "Start to debug.cpuProfile with given path",
	}
	EnableGoTrace = cli.StringFlag{
		Name:  "gotrace",
		Usage: "Start to debug.startGoTrace with given path.",
	}
	OutputPath = cli.StringFlag{
		Name:  "outputpath",
		Usage: "Dir path of output",
	}
	// commands
	ExecuteCommand = cli.Command{
		Name:  "execute",
		Usage: "request transactions to given nodes",
		Subcommands: []cli.Command{
			{
				Name:   "transfer",
				Usage:  "sends only transfer transactions",
				Action: transferTx,
				Flags: []cli.Flag{
					ChainIDFlag,
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
					EnableCpuProfile,
					EnableGoTrace,
				},
			},
		},
	}
)

type doWork func(taskID string, count uint64)

// transfer send transactions given cli context
func transferTx(ctx *cli.Context) error {
	color.Yellow.Println(">>> Setup to send transfer transactions <<<")

	// setup
	cfg, err := parseConfig(ctx)
	if err != nil {
		return err
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		color.Yellow.Println("> failed to marshal config : ", err)
	} else {
		color.Yellow.Println(">", string(b))
	}

	testCtx := parseContext(cfg)
	if testCtx.TxCount < 1 && testCtx.Duration <= 0 {
		utils.Fatalf("invalid tx count %d with duration %s", testCtx.TxCount, testCtx.Duration.String())
	}

	setupContext(testCtx)
	color.Yellow.Println("> Success to setup context")
	defer tearDownContext(testCtx)

	// start to test transfer transactions after initial delay
	if testCtx.InitDelay > 0 {
		time.Sleep(time.Duration(testCtx.InitDelay))
	}
	color.Cyan.Println(">>> Start to send transfer transactions <<<")
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

	// record test results
	result := make(map[string]interface{})
	nodeCtx := testCtx.NodeContexts[0]
	nodeClient, rpcErr := rpc.DialContext(context.Background(), nodeCtx.url)

	// start time, end time, test args
	layout := "15:04:05.999"
	result["StartTime"] = startTime.Format(layout)
	result["EndTime"] = endTime.Format(layout)
	result["Repeat"] = testCtx.TxCount
	result["Duration"] = testCtx.Duration.String()
	result["Interval"] = testCtx.TxInterval

	// node result
	type nodeResult struct {
		Url     string
		Request uint
		Success uint
		Fail    uint
	}
	var nodeResults []nodeResult
	for _, nodeCtx := range testCtx.NodeContexts {
		nodeResults = append(nodeResults, nodeResult{
			Url:     nodeCtx.url,
			Request: nodeCtx.summary.try,
			Success: nodeCtx.summary.success,
			Fail:    nodeCtx.summary.fail,
		})
	}
	result["Nodes"] = nodeResults

	// address's transactions result
	type addressResult struct {
		Hex                string
		Request            uint
		Success            uint
		Fail               uint
		FirstTxHash        string
		FirstTxBlockNumber uint64
	}
	var addressResults []addressResult
	for _, addrCtx := range testCtx.AddressContexts {
		// getting first block number including address's transaction
		firstBlock := uint64(0)
		if rpcErr == nil && addrCtx.summary.firstTxHash != "" {
			hash := common.HexToHash(addrCtx.summary.firstTxHash)
			var result map[string]interface{}
			if err := nodeClient.Call(&result, "berith_getTransactionByHash", hash); err == nil {
				if _, ok := result["blockNumber"]; ok {
					firstBlock, _ = hexutil.DecodeUint64(fmt.Sprintf("%v", result["blockNumber"]))
				}
			}
		}
		addressResults = append(addressResults, addressResult{
			Hex:                addrCtx.account.Address.Hex(),
			Request:            addrCtx.summary.try,
			Success:            addrCtx.summary.success,
			Fail:               addrCtx.summary.fail,
			FirstTxHash:        addrCtx.summary.firstTxHash,
			FirstTxBlockNumber: firstBlock,
		})
	}
	result["Addresses"] = addressResults

	// execute template
	tmpl, err := template.New("").Parse(transferResultTemplate)
	if err != nil {
		fmt.Println("failed to execute template after test completion")
		return err
	}
	var out bytes.Buffer
	err = tmpl.Execute(&out, result)
	if err != nil {
		fmt.Println("failed to execute template after test completion")
		return err
	}

	// display to std out
	fmt.Println(out.String())

	// write a result file
	if testCtx.OutputPath != "" {
		path := testCtx.OutputPath
		fi, err := os.Stat(path)
		if err == nil {
			if !fi.IsDir() {
				fmt.Println("output path must be directory path not file. path :", path)
				return nil
			}
		} else if !os.IsNotExist(err) {
			fmt.Printf("failed to get output path's stat %v\n", err)
			return nil
		}
		path = filepath.Join(path, "berithench-"+testCtx.TestId+".out")
		fmt.Println(">> try to write output :", path)
		err = ioutil.WriteFile(path, out.Bytes(), 0644)
		if err != nil {
			fmt.Printf("failed to write output %v\n", err)
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
			fmt.Printf("cannot sign a transaction :%v\n", err)
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
		fmt.Println("failed to setup task. addr :", addrCtx.account.Address.Hex())
	}
}

// setupContext setup test context such as calc nonce.
func setupContext(ctx *berithenchContext) {
	// setup test id
	layout := "150405"
	ctx.TestId = time.Now().Format(layout)

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

	// start profile
	if !ctx.EnableCpuProfile && !ctx.EnableGoTrace {
		return
	}

	for i, nodeCtx := range ctx.NodeContexts {
		nodeClient, err := rpc.DialContext(context.Background(), nodeCtx.url)
		if err != nil {
			utils.Fatalf("cant create a rpc client. url : %s, err : %v", nodeCtx.url, err)
		}

		if ctx.EnableCpuProfile {
			name := "cpu-profile-node" + strconv.Itoa(i) + "-" + ctx.TestId
			err = nodeClient.Call(nil, "debug_startCPUProfile", name)
			if err != nil {
				if strings.Contains(err.Error(), "already in progress") {
					color.Yellow.Printf("> skip cpu profile because already in progress. node %s\n", nodeCtx.url)
					continue
				}
				utils.Fatalf("failed to start cpu profile: %v", err)
			}
			color.Yellow.Printf("> started cpu profile. node %s : %s\n", nodeCtx.url, name)
			nodeCtx.startCpuProfile = true
		}

		if ctx.EnableGoTrace {
			name := "go-trace-node" + strconv.Itoa(i) + "-" + ctx.TestId
			err = nodeClient.Call(nil, "debug_startGoTrace", name)
			if err != nil {
				if strings.Contains(err.Error(), "already in progress") {
					color.Yellow.Printf("> skip cpu profile because already in progress. node %s\n", nodeCtx.url)
					continue
				}
				utils.Fatalf("failed to start cpu profile: %v", err)
			}
			color.Yellow.Printf("> started go trace. node %s : %s\n", nodeCtx.url, name)
			nodeCtx.startGoTrace = true
		}
	}
}

// tearDownContext clean up test context if needed
func tearDownContext(ctx *berithenchContext) {
	// stop profile
	if ctx.EnableCpuProfile || ctx.EnableGoTrace {
		for _, nodeCtx := range ctx.NodeContexts {
			nodeClient, err := rpc.DialContext(context.Background(), nodeCtx.url)
			if err != nil {
				fmt.Println("cannot create a rpc client. url :", nodeCtx.url)
				continue
			}

			if ctx.EnableCpuProfile && nodeCtx.startCpuProfile {
				err = nodeClient.Call(nil, "debug_stopCPUProfile")
				if err != nil {
					fmt.Printf("failed to stop cpu profile: %v\n", err)
				}
			}

			if ctx.EnableGoTrace && nodeCtx.startGoTrace {
				err = nodeClient.Call(nil, "debug_stopGoTrace")
				if err != nil {
					fmt.Printf("failed to stop go trace: %v\n", err)
				}
			}
		}
		return
	}

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

// parseContext parse berithench config to test context with unlock all
func parseContext(config *berithenchConfig) *berithenchContext {
	ctx := berithenchContext{
		ChainID:          big.NewInt(config.ChainID),
		Duration:         parseDuration(config),
		TxCount:          config.TxCount,
		TxInterval:       config.TxInterval,
		InitDelay:        config.InitDelay,
		OutputPath:       config.OutputPath,
		EnableCpuProfile: config.EnableCpuProfile,
		EnableGoTrace:    config.EnableGoTrace,
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
