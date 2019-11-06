package core

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/BerithFoundation/berith-chain/core/vm"

	"github.com/BerithFoundation/berith-chain/accounts/keystore"
	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/core/state"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/params"
)

type txdata struct {
	to       common.Address
	data     []byte
	value    *big.Int
	base     types.JobWallet
	target   types.JobWallet
	nonce    uint64
	gas      uint64
	gasPrice *big.Int
}

var (
	eth, _ = new(big.Int).SetString("1000000000000000000", 10)

	memdb      = state.NewDatabase(berithdb.NewMemDatabase())
	stateDB, _ = state.New(common.Hash{}, memdb)
	header     = types.Header{
		ParentHash:  common.Hash{},
		UncleHash:   common.Hash{},
		Coinbase:    common.Address{},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.BytesToBloom([]byte("")),
		Difficulty:  big.NewInt(1234),
		Number:      big.NewInt(1),
		GasLimit:    200000000000,
		GasUsed:     0,
		Time:        big.NewInt(time.Now().Unix()),
		Extra:       make([]byte, 0),
		MixDigest:   common.Hash{},
		Nonce:       types.EncodeNonce(0),
	}
	ks = keystore.NewKeyStore(os.TempDir()+"/keystore/", keystore.StandardScryptN, keystore.StandardScryptP)

	vmConfig = vm.Config{
		EnablePreimageRecording: false,
		EWASMInterpreter:        "",
		EVMInterpreter:          "",
	}

	datas = []txdata{
		txdata{
			to:       common.Address{}, //if try to stake or unstake it automatically set to senders address
			nonce:    0,
			value:    new(big.Int).Mul(big.NewInt(100000), eth),
			data:     make([]byte, 0),
			base:     types.Main,
			target:   types.Stake,
			gas:      21000,
			gasPrice: big.NewInt(2000),
		},
		txdata{
			to:       common.BytesToAddress([]byte("to")),
			nonce:    0,
			value:    new(big.Int).Mul(big.NewInt(100000), eth),
			data:     make([]byte, 0),
			base:     types.Stake,
			target:   types.Main,
			gas:      21000,
			gasPrice: big.NewInt(2000),
		},
		txdata{
			to:       common.Address{}, //if try to stake or unstake it automatically set to senders address
			nonce:    0,
			value:    new(big.Int).Mul(big.NewInt(100000), eth),
			data:     make([]byte, 0),
			base:     types.Main,
			target:   types.Stake,
			gas:      21000,
			gasPrice: big.NewInt(2000),
		},
	}
)

// func Test01(t *testing.T) {
// 	addr, _ := ks.NewAccount("0000")
// 	println(addr.Address.Hex())
// }

func TestApplyAndRevertTransaction(t *testing.T) {

	sender, err := ks.NewAccount("1234")

	if err != nil {
		t.Error(err)
	}

	defer func() {
		ks.Delete(sender, "1234")
	}()

	for _, data := range datas {

		init, _ := new(big.Int).SetString("100000000000000000000000000000000000000000000000", 10)
		stateDB.AddBalance(sender.Address, init)
		stateDB.AddStakeBalance(sender.Address, init, big.NewInt(1))
		stateDB.Commit(true)

		if data.base == types.Stake || data.target == types.Stake {
			data.to = sender.Address
		}

		tx := types.NewTransaction(data.nonce, data.to, data.value, data.gas, data.gasPrice, data.data, data.base, data.target)
		tx, err = ks.SignTxWithPassphrase(sender, "1234", tx, params.TestnetChainConfig.ChainID)
		if err != nil {
			t.Error(err)
		}
		msg, err := tx.AsMessage(types.NewEIP155Signer(params.TestnetChainConfig.ChainID))
		if err != nil {
			t.Error(err)
		}
		author := common.BytesToAddress([]byte("gas"))
		ctx := NewEVMContext(msg, &header, nil, &author)
		gp := new(GasPool)
		gp.AddGas(header.GasLimit)
		evm := vm.NewEVM(ctx, stateDB, params.TestnetChainConfig, vmConfig)

		snap := stateDB.Snapshot()

		stateDB.AddPenalty(sender.Address, header.Number)

		pen := stateDB.GetPenalty(sender.Address)
		penUdt := stateDB.GetPenaltyUpdated(sender.Address)

		if pen != 1 || penUdt.Cmp(header.Number) != 0 {
			t.Errorf("expected result is [1,%s] but [%d,%s]", header.Number.String(), pen, penUdt.String())
		}

		stateDB.RemovePenalty(sender.Address, header.Number)

		pen = stateDB.GetPenalty(sender.Address)
		penUdt = stateDB.GetPenaltyUpdated(sender.Address)

		if pen != 0 || penUdt.Cmp(header.Number) != 0 {
			t.Errorf("expected result is [1,%s] but [%d,%s]", header.Number.String(), pen, penUdt.String())
		}

		originFrom := getBalances(sender.Address, stateDB)
		originTo := getBalances(*tx.To(), stateDB)

		gasAmt := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))

		exptFrom := getBalances(sender.Address, stateDB)
		exptTo := getBalances(sender.Address, stateDB)

		_, _, _, err = ApplyMessage(evm, msg, gp)
		if err != nil {
			t.Error(err)
		}

		exptFrom.bal = new(big.Int).Sub(exptFrom.bal, gasAmt)

		if tx.Target() == types.Main && tx.Base() == types.Main {
			exptFrom.bal = new(big.Int).Sub(exptFrom.bal, tx.Value())
			exptTo.bal = new(big.Int).Add(exptTo.bal, tx.Value())
		} else if tx.Target() == types.Stake {
			exptFrom.bal = new(big.Int).Sub(exptFrom.bal, tx.Value())
			exptFrom.stk = new(big.Int).Add(exptFrom.stk, tx.Value())
			exptFrom.stkUdt = new(big.Int).Set(header.Number)
			exptTo = exptFrom

		} else {
			exptFrom.bal = new(big.Int).Add(exptFrom.bal, exptFrom.stk)
			exptFrom.stk = big.NewInt(0)
			exptFrom.stkUdt = new(big.Int).Set(header.Number)
			exptTo = exptFrom
		}

		resultFrom := getBalances(msg.From(), stateDB)
		resultTo := getBalances(*msg.To(), stateDB)

		if err := exptFrom.Compare(resultFrom); err != nil {
			t.Error(err)
		}
		if err := exptTo.Compare(resultTo); err != nil {
			t.Error(err)
		}

		if err := checkBigInt(stateDB.GetBalance(author), gasAmt); err != nil {
			t.Error(err)
		}

		stateDB.RevertToSnapshot(snap)

		revertedFrom := getBalances(msg.From(), stateDB)
		revertedTo := getBalances(*msg.To(), stateDB)

		if err := originFrom.Compare(revertedFrom); err != nil {
			t.Error(err)
		}
		if err := originTo.Compare(revertedTo); err != nil {
			t.Error(err)
		}
		if err := checkBigInt(stateDB.GetBalance(author), big.NewInt(0)); err != nil {
			t.Error(err)
		}
	}
}

type balances struct {
	bal    *big.Int
	stk    *big.Int
	stkUdt *big.Int
}

func (b balances) Compare(other balances) error {
	if err := checkBigInt(b.bal, other.bal); err != nil {
		return err
	}
	if err := checkBigInt(b.stk, other.stk); err != nil {
		return err
	}
	if err := checkBigInt(b.stkUdt, other.stkUdt); err != nil {
		return err
	}
	return nil
}

func getBalances(addr common.Address, stateDB *state.StateDB) balances {
	return balances{
		bal:    new(big.Int).Set(stateDB.GetBalance(addr)),
		stk:    new(big.Int).Set(stateDB.GetStakeBalance(addr)),
		stkUdt: new(big.Int).Set(stateDB.GetStakeUpdated(addr)),
	}
}

func checkBigInt(a, b *big.Int) error {
	if a.Cmp(b) != 0 {
		return fmt.Errorf("expected %s but, %s", a.String(), b.String())
	}
	return nil
}
