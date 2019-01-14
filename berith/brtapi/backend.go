package brtapi

import (
	"context"
	"math/big"

	"bitbucket.org/ibizsoftware/berith-chain/accounts"
	"bitbucket.org/ibizsoftware/berith-chain/common"
	"bitbucket.org/ibizsoftware/berith-chain/core/types"
	"bitbucket.org/ibizsoftware/berith-chain/params"
	"bitbucket.org/ibizsoftware/berith-chain/rpc"
)

//Backend backend of berith service
type Backend interface {
	AccountManager() *accounts.Manager
	SuggestPrice(ctx context.Context) (*big.Int, error)
	GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error)

	SendTx(ctx context.Context, signedTx *types.Transaction) error

	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block
}

//GetAPIs get apis of berith serivce
func GetAPIs(b Backend) []rpc.API {
	nonceLock := new(AddrLocker)

	return []rpc.API{
		{
			Namespace: "berith",
			Version:   "1.0",
			Service:   NewPrivateBerithAPI(b, nonceLock),
			Public:    false,
		},
	}
}
