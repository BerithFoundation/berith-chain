// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package berith

import (
	"context"
	"errors"
	"math/big"

	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/berith/downloader"
	"github.com/BerithFoundation/berith-chain/berith/gasprice"
	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/common/math"
	"github.com/BerithFoundation/berith-chain/core"
	"github.com/BerithFoundation/berith-chain/core/bloombits"
	"github.com/BerithFoundation/berith-chain/core/state"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/core/vm"
	"github.com/BerithFoundation/berith-chain/event"
	"github.com/BerithFoundation/berith-chain/params"
	"github.com/BerithFoundation/berith-chain/rpc"
)

// BerAPIBackend implements berithapi.Backend for full nodes
type BerAPIBackend struct {
	e   *Berith
	gpo *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *BerAPIBackend) ChainConfig() *params.ChainConfig {
	return b.e.chainConfig
}

func (b *BerAPIBackend) CurrentBlock() *types.Block {
	return b.e.blockchain.CurrentBlock()
}

func (b *BerAPIBackend) SetHead(number uint64) {
	b.e.protocolManager.downloader.Cancel()
	b.e.blockchain.SetHead(number)
}

func (b *BerAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.e.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.e.blockchain.CurrentBlock().Header(), nil
	}
	return b.e.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *BerAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.e.blockchain.GetHeaderByHash(hash), nil
}

func (b *BerAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.e.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.e.blockchain.CurrentBlock(), nil
	}
	return b.e.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *BerAPIBackend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.BlockByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.e.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.e.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		block := b.e.blockchain.GetBlock(hash, header.Number.Uint64())
		if block == nil {
			return nil, errors.New("header found, but block body is missing")
		}
		return block, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BerAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.e.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.e.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *BerAPIBackend) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.StateAndHeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := b.HeaderByHash(ctx, hash)
		if err != nil {
			return nil, nil, err
		}
		if header == nil {
			return nil, nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.e.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, nil, errors.New("hash is not currently canonical")
		}
		stateDb, err := b.e.BlockChain().StateAt(header.Root)
		return stateDb, header, err
	}
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *BerAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.e.blockchain.GetBlockByHash(hash), nil
}

func (b *BerAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return b.e.blockchain.GetReceiptsByHash(hash), nil
}

func (b *BerAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	receipts := b.e.blockchain.GetReceiptsByHash(hash)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *BerAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.e.blockchain.GetTdByHash(blockHash)
}

func (b *BerAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.e.BlockChain(), nil)
	return vm.NewEVM(context, state, b.e.chainConfig, *b.e.blockchain.GetVMConfig()), vmError, nil
}

func (b *BerAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.e.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *BerAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.e.BlockChain().SubscribeChainEvent(ch)
}

func (b *BerAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.e.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *BerAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.e.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *BerAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.e.BlockChain().SubscribeLogsEvent(ch)
}

func (b *BerAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.e.txPool.AddLocal(signedTx)
}

func (b *BerAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.e.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *BerAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.e.txPool.Get(hash)
}

func (b *BerAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.e.txPool.State().GetNonce(addr), nil
}

func (b *BerAPIBackend) Stats() (pending int, queued int) {
	return b.e.txPool.Stats()
}

func (b *BerAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.e.TxPool().Content()
}

func (b *BerAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.e.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *BerAPIBackend) Downloader() *downloader.Downloader {
	return b.e.Downloader()
}

func (b *BerAPIBackend) ProtocolVersion() int {
	return b.e.BerVersion()
}

func (b *BerAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *BerAPIBackend) ChainDb() berithdb.Database {
	return b.e.ChainDb()
}

func (b *BerAPIBackend) EventMux() *event.TypeMux {
	return b.e.EventMux()
}

func (b *BerAPIBackend) AccountManager() *accounts.Manager {
	return b.e.AccountManager()
}

func (b *BerAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.e.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *BerAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.e.bloomRequests)
	}
}
