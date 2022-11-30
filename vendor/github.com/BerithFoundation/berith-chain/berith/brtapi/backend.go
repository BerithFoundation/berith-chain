package brtapi

import (
	"github.com/BerithFoundation/berith-chain/core/state"
	"github.com/BerithFoundation/berith-chain/miner"
	"context"
	"math/big"

	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/params"
	"github.com/BerithFoundation/berith-chain/rpc"
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
