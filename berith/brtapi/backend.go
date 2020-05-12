package brtapi

import (
	"context"
	"berith-chain/core/state"
	"berith-chain/miner"
	"math/big"

	"berith-chain/accounts"
	"berith-chain/common"
	"berith-chain/core/types"
	"berith-chain/params"
	"berith-chain/rpc"
)

//Backend backend of berith service
type Backend interface {
	AccountManager() *accounts.Manager
	SuggestPrice(ctx context.Context) (*big.Int, error)
	GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error)

	SendTx(ctx context.Context, signedTx *types.Transaction) error

	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block

	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error)
}

//GetAPIs get apis of berith serivce
func GetAPIs(b Backend, miner *miner.Miner) []rpc.API {
	nonceLock := new(AddrLocker)

	return []rpc.API{
		{
			Namespace: "berith",
			Version:   "1.0",
			Service:   NewPrivateBerithAPI(b, miner, nonceLock),
			Public:    false,
		},
	}
}
