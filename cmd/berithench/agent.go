package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/common/hexutil"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/rlp"
	"github.com/BerithFoundation/berith-chain/rpc"
	"github.com/gookit/color"

	"github.com/BerithFoundation/berith-chain/accounts/keystore"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	transfer = 1 + iota
	stake
	unstake
	contractCall
)

var (
	TxValueFlag = cli.Uint64Flag{
		Name:  "txvalue",
		Usage: "Value of transaciton",
	}

	StakeValueFlag = cli.Uint64Flag{
		Name:  "stakevalue",
		Usage: "Value of stake",
	}

	ContractValueFlag = cli.Uint64Flag{
		Name:  "contractvalue",
		Usage: "Value of contract call",
	}

	GasLimitFlag = cli.Uint64Flag{
		Name:  "gaslimit",
		Usage: "Gas limit",
	}

	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "password of agent accounts",
	}

	GasPriceFlag = cli.Uint64Flag{
		Name:  "gasprice",
		Usage: "Gas price",
	}

	StakeCountFlag = cli.StringFlag{
		Name:  "stakecount",
		Usage: "Amount of transaction for staking",
	}

	ContractCountFlag = cli.StringFlag{
		Name:  "contractcount",
		Usage: "Amount of transaction for calling contract",
	}

	AgentCommand = cli.Command{
		Name:   "agent",
		Usage:  "run a agent to transfer token automatically",
		Action: runAgent,
		Flags: []cli.Flag{
			ChainIDFlag,
			NodesFlag,
			ConfigFileFlag,
			KeystoreFlag,
			PasswordFlag,
			TxCountFlag,
			TxIntervalFlag,
			InitDelay,
			StakeCountFlag,
			ContractCountFlag,
			TxValueFlag,
			StakeValueFlag,
			ContractValueFlag,
			GasPriceFlag,
			GasLimitFlag,
		},
		Subcommands: []cli.Command{
			AgentAccountsCommand,
		},
	}
	AgentAccountsCommand = cli.Command{
		Name:   "accounts",
		Usage:  "manage accounts of agent",
		Action: exportAccount,
		Flags: []cli.Flag{
			KeystoreFlag,
			PasswordFlag,
			ConfigFileFlag,
		},
	}
)

type node struct {
	url   string
	block *big.Int
}

type txResult struct {
	tx   common.Hash
	from int
}

type txsResult struct {
	txs  []common.Hash
	from string
}

type contract struct {
	address common.Address
	data    string
}

type agent struct {
	cfg         *berithenchConfig
	stopCh      chan bool
	newBlockCh  chan *types.Block
	ticker      *time.Ticker
	blockTicker *time.Ticker
	txSub       []chan bool
	blockSub    []chan bool
	keystore    *keystore.KeyStore
	nodes       map[string]node
	logCh       chan string
	errCh       chan error
	txCh        chan txResult
	txsCh       chan txsResult
	txMap       map[common.Hash]bool
	contracts   []contract
}

type blockResult struct {
	Transactions []common.Hash `json:"transactions"`
}

func exportAccount(ctx *cli.Context) error {

	cfg, err := parseConfig(ctx)

	if err != nil {
		return err
	}

	println(cfg.Keystore)
	println(cfg.Password)

	if len(cfg.Keystore) == 0 {
		return errors.New("invalid keystore directory")
	}

	if len(cfg.Password) == 0 {
		return errors.New("invalid password")
	}

	ks := keystore.NewKeyStore(cfg.Keystore, keystore.StandardScryptN, keystore.StandardScryptP)
	if ks == nil {
		return errors.New("invalid keystore directory")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		println()
		println(" 1.create new accounts")
		println(" 2.export accounts")
		println(" others.exit")
		print("select your menu : ")
		input, _, _ := reader.ReadLine()

		switch string(input) {
		case "1":
			println()
			print("how many accounts do you create[1 ~ n, others.exit] : ")
			input, _, _ := reader.ReadLine()

			n, err := strconv.Atoi(string(input))
			if err != nil || n <= 0 {
				break
			}

			created := 0
			for i := len(ks.Accounts()); i < n; i++ {
				if acc, err := ks.NewAccount(cfg.Password); err == nil {
					created++
					println(acc.Address.Hex())
				}
			}
			println()
			println(created, "accounts is created successfully")
			break
		case "2":
			println()
			print("what teplate do you taken[1.json, 2.alloc, 3.source code, others.exit] : ")
			input, _, _ := reader.ReadLine()

			switch string(input) {
			case "1":
				result := make([]common.Address, 0)
				for _, acc := range ks.Accounts() {
					result = append(result, acc.Address)
				}

				if jsondata, err := json.MarshalIndent(result, "", "\t"); err == nil {
					println(string(jsondata))
				}
			case "2":

				type alloc struct {
					Balance string `json:"balance"`
				}

				result := make(map[string]alloc, 0)
				println()
				print("how many provide token to each address[1 ~ n(eth), others.exit]: ")
				input, _, _ := reader.ReadLine()

				token, err := strconv.Atoi(string(input))

				if err != nil || token < 0 {
					break
				}

				balance := new(big.Int).Mul(big.NewInt(int64(token)), big.NewInt(1e+8))

				for _, acc := range ks.Accounts() {
					result[acc.Address.Hex()[2:]] = alloc{Balance: balance.String()}
				}

				if jsondata, err := json.MarshalIndent(result, "", "\t"); err == nil {
					println(string(jsondata))
				}
			case "3":

				println()
				print("how many provide token to each address[1 ~ n(eth), others.exit]: ")
				input, _, _ := reader.ReadLine()

				token, err := strconv.Atoi(string(input))

				if err != nil || token < 0 {
					break
				}

				balance := new(big.Int).Mul(big.NewInt(int64(token)), big.NewInt(1e+8))

				for _, acc := range ks.Accounts() {
					fmt.Printf("common.HexToAddress(\"%s\"): {Balance: common.StringToBig(\"%s\")},\n", acc.Address.Hex(), balance.String())
				}

			}
		default:
			return nil
		}
	}
}

func newAgent(cfg *berithenchConfig) (*agent, error) {

	configs, _ := json.MarshalIndent(cfg, "", "\t")

	println("===============[CONFIG]==============")
	println(string(configs))

	agent := agent{
		cfg:         cfg,
		keystore:    keystore.NewKeyStore(cfg.Keystore, keystore.StandardScryptN, keystore.StandardScryptP),
		stopCh:      make(chan bool),
		newBlockCh:  make(chan *types.Block),
		nodes:       make(map[string]node),
		ticker:      time.NewTicker(time.Millisecond * time.Duration(cfg.TxInterval) / time.Duration(cfg.TxCount+cfg.StakeCount+cfg.ContractCount)),
		blockTicker: time.NewTicker(time.Second),
		txSub:       make([]chan bool, 0),
		blockSub:    make([]chan bool, 0),
		logCh:       make(chan string),
		errCh:       make(chan error),
		txCh:        make(chan txResult),
		txsCh:       make(chan txsResult),
		txMap:       make(map[common.Hash]bool),
	}

	for keyAmt := uint64(len(agent.keystore.Accounts())); keyAmt <= cfg.TxCount+cfg.StakeCount+cfg.ContractCount; keyAmt++ {
		//agent.keystore.NewAccount(cfg.Password)
		return nil, errors.New("not enough account to run agent")
	}

	for _, url := range cfg.Nodes {
		agent.nodes[url] = node{
			url:   url,
			block: big.NewInt(0),
		}
	}

	if agent.cfg.ContractCount > 0 {
		contracts := make([]contract, 0)
		for addr, data := range cfg.ContractData {
			for _, datum := range data {
				contracts = append(contracts, contract{
					address: common.HexToAddress(addr),
					data:    datum,
				})
			}
		}
		if len(contracts) == 0 {
			return nil, errors.New("no available contract data given")
		}
		agent.contracts = contracts
	}

	return &agent, nil
}

func (agent *agent) run() {
	go agent.logLoop()
	go agent.txLoop()
	for url := range agent.nodes {
		ch := make(chan bool)
		go agent.blockLoop(url, ch)
		agent.blockSub = append(agent.blockSub, ch)
	}
	for i := 0; i < int(agent.cfg.TxCount+agent.cfg.StakeCount+agent.cfg.ContractCount); i++ {
		ch := make(chan bool)
		agent.keystore.Unlock(agent.keystore.Accounts()[i], agent.cfg.Password)
		go agent.transferLoop(agent.getNodeByIdx(i), i, ch)
		agent.txSub = append(agent.txSub, ch)
	}
	agent.tickerLoop()

}

func (agent *agent) logLoop() {
	for {
		select {
		case log := <-agent.logCh:
			println(log)
		case err := <-agent.errCh:
			color.Red.Println(err.Error())
		case <-agent.stopCh:
			return
		}
	}
}

func (agent *agent) txLoop() {
	for {
		select {
		case txs := <-agent.txsCh:
			txsFromAgent := 0
			for _, tx := range txs.txs {
				if _, ok := agent.txMap[tx]; ok {
					txsFromAgent++
					delete(agent.txMap, tx)
				}
			}
			agent.logCh <- fmt.Sprintf("[%s]%d Txs is processed, left %d txs", txs.from, txsFromAgent, len(agent.txMap))
		case tx := <-agent.txCh:
			agent.txMap[tx.tx] = true
			agent.logCh <- fmt.Sprintf("[Acc%d]Imported Tx \"%s\"", tx.from, tx.tx.Hex())
		case <-agent.stopCh:
			return
		}
	}
}

func (agent *agent) blockLoop(url string, ch chan bool) {
	client, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		agent.errCh <- err
		return
	}
	for {
		select {
		case <-ch:
			result := ""
			if err := client.CallContext(context.Background(), &result, "berith_blockNumber", []string{}); err != nil {
				agent.errCh <- err
				continue
			}

			blockNumber, ok := new(big.Int).SetString(result[2:], 16)
			if !ok {
				continue
			}
			if blockNumber.Cmp(agent.nodes[url].block) == 1 {
				agent.nodes[url].block.Set(blockNumber)
				target := hexutil.EncodeBig(agent.nodes[url].block)
				holder := blockResult{
					Transactions: make([]common.Hash, 0),
				}
				if err := client.CallContext(context.Background(), &holder, "berith_getBlockByNumber", target, false); err != nil {
					agent.errCh <- err
				}
				if holder.Transactions != nil && len(holder.Transactions) > 0 {
					agent.txsCh <- txsResult{
						txs:  holder.Transactions,
						from: url,
					}
				}

			}
		case <-agent.stopCh:
			return
		}
	}
}

func (agent *agent) tickerLoop() {
	index := 0
	for {
		select {
		case <-agent.blockTicker.C:
			for _, ch := range agent.blockSub {
				ch <- true
			}
		case <-agent.ticker.C:
			agent.txSub[index] <- true
			index = (index + 1) % len(agent.txSub)
		}
	}
}

func (agent *agent) transferLoop(url string, index int, ch chan bool) {
	var nonce uint64

	txType := transfer

	if index >= int(agent.cfg.TxCount) {
		txType = stake
	}

	if index >= int(agent.cfg.TxCount+agent.cfg.StakeCount) {
		txType = contractCall
	}

	client, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		agent.errCh <- err
	}

	var nonceStr string

	err = client.CallContext(context.Background(), &nonceStr, "berith_getTransactionCount", agent.keystore.Accounts()[index].Address, "pending")

	if err != nil {
		agent.errCh <- err
		return
	}
	nonce, err = hexutil.DecodeUint64(nonceStr)

	if err != nil {
		agent.errCh <- err
		return
	}

	agent.logCh <- fmt.Sprintf("[Acc%d] is Running", index)
	for {
		select {
		case <-ch:
			var base, target types.JobWallet
			base = types.Main
			target = types.Main
			if txType == stake {
				base = types.Main
				target = types.Stake
			}
			if txType == unstake {
				base = types.Stake
				target = types.Main
			}

			//[TODO] : change base and target in case of transaction type is steak or unstake

			accs := agent.keystore.Accounts()

			from := accs[index]

			to := accs[uint64(index+1)%agent.cfg.TxCount].Address
			if txType == stake || txType == unstake {
				to = accs[index].Address
			}
			if txType == contractCall {
				to = agent.contracts[(index-int(agent.cfg.TxCount+agent.cfg.StakeCount))%len(agent.contracts)].address
			}

			value := big.NewInt(int64(agent.cfg.TxValue))
			if txType == stake || txType == unstake {
				value = big.NewInt(int64(agent.cfg.StakeValue))
			}
			if txType == contractCall {
				value = big.NewInt(int64(agent.cfg.ContractValue))
			}
			value.Mul(value, big.NewInt(1e+8))

			txData := make([]byte, 0)
			if txType == contractCall {
				txData = common.Hex2Bytes(agent.contracts[(index-int(agent.cfg.TxCount+agent.cfg.StakeCount))%len(agent.contracts)].data)
			}

			gasLimit := agent.cfg.GasLimit
			gasPrice := big.NewInt(int64(agent.cfg.GasPrice))

			tx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, txData, base, target)

			tx, err := agent.keystore.SignTx(from, tx, big.NewInt(agent.cfg.ChainID))

			if err != nil {
				agent.errCh <- err
				return
			}

			data, err := rlp.EncodeToBytes(tx)
			if err != nil {
				agent.errCh <- err
				return
			}

			resp := common.Hash{}

			if err := client.CallContext(context.Background(), &resp, "berith_sendRawTransaction", common.ToHex(data)); err != nil {
				agent.errCh <- err
				continue
			}

			emptyHash := common.Hash{}
			if resp != emptyHash {
				agent.txCh <- txResult{
					tx:   resp,
					from: index,
				}
			}

			//println("RESULT : ", resp.Hex())
			if txType == stake {
				txType = unstake
			} else if txType == unstake {
				txType = stake
			}
			nonce++
		case <-agent.stopCh:
			return
		}
	}

}

func (agent *agent) getNodeByIdx(index int) (node string) {
	i := index % len(agent.nodes)
	n := 0
	for node = range agent.nodes {
		if i == n {
			return
		}
		n++
	}
	return
}

func runAgent(ctx *cli.Context) error {
	config, err := parseConfig(ctx)

	if err != nil {
		return err
	}

	agent, err := newAgent(config)

	if err != nil {
		return err
	}

	go agent.run()

	<-agent.stopCh

	return nil
}
