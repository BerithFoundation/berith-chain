package core

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/BerithFoundation/berith-chain/core/vm"
	"github.com/BerithFoundation/berith-chain/crypto/secp256k1"

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

	pvk := new(ecdsa.PrivateKey)

	pvk.PublicKey.Curve = secp256k1.S256()
	pvk.D = big.NewInt(3360)
	pvk.PublicKey.X, pvk.PublicKey.Y = secp256k1.S256().ScalarBaseMult(big.NewInt(3360).Bytes())

	from := common.HexToAddress("Bx810722274468C2E5dEE8Aabd41aE61fA4d1A5cDa")

	signer := types.NewEIP155Signer(big.NewInt(206))

	for _, data := range datas {

		init, _ := new(big.Int).SetString("100000000000000000000000000000000000000000000000", 10)
		stateDB.AddBalance(from, init)
		stateDB.AddStakeBalance(from, init, big.NewInt(1))
		stateDB.Commit(true)

		if data.base == types.Stake || data.target == types.Stake {
			data.to = from
		}

		tx := types.NewTransaction(data.nonce, data.to, data.value, data.gas, data.gasPrice, data.data, data.base, data.target)

		tx, err := types.SignTx(tx, signer, pvk)

		if err != nil {
			t.Error(err)
		}

		msg, err := tx.AsMessage(types.NewEIP155Signer(params.TestnetChainConfig.ChainID))
		if err != nil {
			t.Error(err)
		}
		author := common.BytesToAddress([]byte("gas"))
		isBIP4 := params.TestnetChainConfig.IsBIP4(header.Number)
		ctx := NewEVMContext(msg, &header, nil, &author, isBIP4)
		gp := new(GasPool)
		gp.AddGas(header.GasLimit)
		evm := vm.NewEVM(ctx, stateDB, params.TestnetChainConfig, vmConfig)

		snap := stateDB.Snapshot()

		stateDB.AddPenalty(from, header.Number)

		pen := stateDB.GetPenalty(from)
		penUdt := stateDB.GetPenaltyUpdated(from)

		if pen != 1 || penUdt.Cmp(header.Number) != 0 {
			t.Errorf("expected result is [1,%s] but [%d,%s]", header.Number.String(), pen, penUdt.String())
		}

		stateDB.RemovePenalty(from, header.Number)

		pen = stateDB.GetPenalty(from)
		penUdt = stateDB.GetPenaltyUpdated(from)

		if pen != 0 || penUdt.Cmp(header.Number) != 0 {
			t.Errorf("expected result is [1,%s] but [%d,%s]", header.Number.String(), pen, penUdt.String())
		}

		originFrom := getBalances(from, stateDB)
		originTo := getBalances(*tx.To(), stateDB)

		gasAmt := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))

		exptFrom := getBalances(from, stateDB)
		exptTo := getBalances(from, stateDB)

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

func TestGetCanTransferFunc(t *testing.T) {
	type testData struct {
		isBIP4 bool
		target types.JobWallet
		want vm.CanTransferFunc
	}

	tests := []testData {
		testData{isBIP4: false, target: types.Main, want: CanTransfer},
		testData{isBIP4: false, target: types.Stake, want: CanTransfer},
		testData{isBIP4: true, target: types.Main, want: CanTransfer},
		testData{isBIP4: true, target: types.Stake, want: CanTransferBIP4},
	}

	for _, test := range tests {
		resultFuncName := runtime.FuncForPC(reflect.ValueOf(getCanTransferFunc(test.isBIP4, test.target)).Pointer()).Name()
		expectedFuncName := runtime.FuncForPC(reflect.ValueOf(test.want).Pointer()).Name()

		if resultFuncName != expectedFuncName {
			t.Errorf("expected %s but %s", expectedFuncName, resultFuncName)
		}
	}
}

func TestCheckStakeBalanceAmount(t *testing.T) {
	type testData struct {
		totalStakingAmount *big.Int
		maximum *big.Int
		want bool
	}

	tests := []testData {
		testData{totalStakingAmount: big.NewInt(49999999), maximum: big.NewInt(50000000), want: true},
		testData{totalStakingAmount: big.NewInt(50000000), maximum: big.NewInt(50000000), want: true},
		testData{totalStakingAmount: big.NewInt(50000001), maximum: big.NewInt(50000000), want: false},
	}

	for _, test := range tests {
		if CheckStakeBalanceAmount(test.totalStakingAmount, test.maximum) != test.want {
			t.Errorf("expected but")
		}
	}
}
