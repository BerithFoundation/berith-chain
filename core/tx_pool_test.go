package core

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"

	"berith-chain/berith/staking"
	"berith-chain/berithdb"
	"berith-chain/common/hexutil"
	"berith-chain/consensus/bsrr"
	"berith-chain/crypto/secp256k1"
	"berith-chain/params"

	"berith-chain/common"
	"berith-chain/core/types"
	"berith-chain/core/vm"
)

func TestTransactionValidate(t *testing.T) {
	pvk := new(ecdsa.PrivateKey)

	pvk.PublicKey.Curve = secp256k1.S256()
	pvk.D = big.NewInt(3360)
	pvk.PublicKey.X, pvk.PublicKey.Y = secp256k1.S256().ScalarBaseMult(big.NewInt(3360).Bytes())

	txs := []txdata{
		txdata{
			to:       common.BytesToAddress([]byte("to")),
			value:    big.NewInt(100000),
			gas:      21000,
			gasPrice: big.NewInt(100000),
			data:     make([]byte, 0),
			base:     4,
			target:   types.Main,
			nonce:    0,
		},
		txdata{
			to:       common.BytesToAddress([]byte("to")),
			value:    big.NewInt(100000),
			gas:      21000,
			gasPrice: big.NewInt(100000),
			data:     make([]byte, 0),
			base:     types.Main,
			target:   4,
			nonce:    0,
		},
		txdata{
			to:       common.BytesToAddress([]byte("to")),
			value:    big.NewInt(100000),
			gas:      21000,
			gasPrice: big.NewInt(100000),
			data:     make([]byte, 0),
			base:     types.Stake,
			target:   types.Stake,
			nonce:    0,
		},
		txdata{
			to:       common.BytesToAddress([]byte("to")),
			value:    big.NewInt(100000),
			gas:      21000,
			gasPrice: big.NewInt(100000),
			data:     make([]byte, 0),
			base:     types.Main,
			target:   types.Stake,
			nonce:    0,
		},
		txdata{
			to:       common.BytesToAddress([]byte("to")),
			value:    big.NewInt(100000),
			gas:      21000,
			gasPrice: big.NewInt(100000),
			data:     make([]byte, 0),
			base:     types.Stake,
			target:   types.Main,
			nonce:    0,
		},
	}

	results := []error{
		types.ErrInvalidJobWallet,
		types.ErrInvalidJobWallet,
		types.ErrInvalidJobWallet,
		ErrInvalidStakeReceiver,
		ErrInvalidStakeReceiver,
	}

	genesis := &Genesis{
		Config:     params.MainnetChainConfig,
		Nonce:      0x00,
		Timestamp:  0x00,
		ExtraData:  hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000810722274468C2E5dEE8Aabd41aE61fA4d1A5cDa0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   94000000,
		Difficulty: big.NewInt(1),
		Mixhash:    common.BytesToHash(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000")),
		Coinbase:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Alloc: map[common.Address]GenesisAccount{
			common.HexToAddress("Bx810722274468C2E5dEE8Aabd41aE61fA4d1A5cDa"): {Balance: common.StringToBig("10000000000000000000000000000")},
		},
	}

	stkDB := new(staking.StakingDB)

	if err := stkDB.CreateDB(os.TempDir()+"/stakingdb/", staking.NewStakers); err != nil {
		t.Error(err)
	}

	memDB := berithdb.NewMemDatabase()

	SetupGenesisBlockWithOverride(memDB, genesis, big.NewInt(0))

	engine := bsrr.NewCliqueWithStakingDB(stkDB, params.TestnetChainConfig.Bsrr, memDB)
	chain, err := NewBlockChain(stkDB, memDB, nil, params.TestnetChainConfig, engine, vm.Config{}, nil)

	if err != nil {
		t.Error(err)
	}

	block := chain.GetBlockByNumber(0)

	println(block.Hash().Hex())

	if err != nil {
		t.Error(err)
	}

	pool := NewTxPool(DefaultTxPoolConfig, params.TestnetChainConfig, chain)

	signer := types.NewEIP155Signer(big.NewInt(206))

	for i, tx := range txs {

		transaction := types.NewTransaction(tx.nonce, tx.to, tx.value, tx.gas, tx.gasPrice, tx.data, tx.base, tx.target)
		transaction, err = types.SignTx(transaction, signer, pvk)

		if err != nil {
			t.Error(err)
		}

		err = pool.AddLocal(transaction)
		if err != results[i] {
			t.Error(err)
		}
	}
}
