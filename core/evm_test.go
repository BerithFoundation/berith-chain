package core

import (
	"math/big"
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
	ks = keystore.NewKeyStore("C:/Users/ibizsoftware/test/keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	vmConfig = vm.Config{
		EnablePreimageRecording: false,
		EWASMInterpreter:        "",
		EVMInterpreter:          "",
	}

	data = txdata{
		to:       ks.Accounts()[0].Address,
		nonce:    0,
		value:    new(big.Int).Mul(big.NewInt(100000), eth),
		data:     make([]byte, 0),
		base:     types.Main,
		target:   types.Stake,
		gas:      21000,
		gasPrice: big.NewInt(2000),
	}
)

// func Test01(t *testing.T) {
// 	addr, _ := ks.NewAccount("0000")
// 	println(addr.Address.Hex())
// }

func Test02(t *testing.T) {

	from := ks.Accounts()[0]

	init, _ := new(big.Int).SetString("100000000000000000000000000000000000000000000000", 10)
	stateDB.AddBalance(from.Address, init)
	stateDB.Commit(true)

	for i, acc := range ks.Accounts() {
		println(i, "=>", acc.Address.Hex())
	}

	tx := types.NewTransaction(data.nonce, data.to, data.value, data.gas, data.gasPrice, data.data, data.base, data.target)
	tx, err := ks.SignTxWithPassphrase(ks.Accounts()[0], "1234", tx, params.TestnetChainConfig.ChainID)
	if err != nil {
		t.Error(err)
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(params.TestnetChainConfig.ChainID))
	if err != nil {
		t.Error(err)
	}
	author := common.BytesToAddress([]byte("gas"))
	println("AUTHOR ===>>>", author.Hex())
	ctx := NewEVMContext(msg, &header, nil, &author)
	gp := new(GasPool)
	gp.AddGas(header.GasLimit)
	evm := vm.NewEVM(ctx, stateDB, params.TestnetChainConfig, vmConfig)

	snap := stateDB.Snapshot()

	println("BLOCKNUMBER ==>>", evm.BlockNumber.String())

	stateDB.AddPenalty(from.Address, header.Number)

	printAccount(stateDB, from.Address)
	printAccount(stateDB, data.to)
	println("GAS ==>>", stateDB.GetBalance(author).String())
	stateDB.RemovePenalty(from.Address, header.Number)

	_, _, _, err = ApplyMessage(evm, msg, gp)
	if err != nil {
		stateDB.RevertToSnapshot(snap)
		println("==================REVERTED================")
	}

	printAccount(stateDB, from.Address)
	printAccount(stateDB, data.to)
	println("GAS ==>>", stateDB.GetBalance(author).String())

	stateDB.RevertToSnapshot(snap)
	println("==================REVERTED================")

	printAccount(stateDB, from.Address)
	printAccount(stateDB, data.to)
	println("GAS ==>>", stateDB.GetBalance(author).String())
}

func printAccount(stateDB *state.StateDB, addr common.Address) {
	println("=================[", addr.Hex(), "]=================")
	println("MAIN ==>>", stateDB.GetBalance(addr).String())
	println("STAKE ==>> [", stateDB.GetStakeBalance(addr).String(), ",", stateDB.GetStakeUpdated(addr).String(), "]")
	println("PENALTY ==>> [", stateDB.GetPenalty(addr), ",", stateDB.GetPenaltyUpdated(addr).String(), "]")
	println()
}
