package brtapi

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
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
