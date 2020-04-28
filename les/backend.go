// Copyright 2016 The go-ethereum Authors
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

// Package les implements the Light Ethereum Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/berith"
	"github.com/BerithFoundation/berith-chain/berith/downloader"
	"github.com/BerithFoundation/berith-chain/berith/filters"
	"github.com/BerithFoundation/berith-chain/berith/gasprice"
	"github.com/BerithFoundation/berith-chain/berith/staking"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/common/hexutil"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/core"
	"github.com/BerithFoundation/berith-chain/core/bloombits"
	"github.com/BerithFoundation/berith-chain/core/rawdb"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/event"
	"github.com/BerithFoundation/berith-chain/internal/berithapi"
	"github.com/BerithFoundation/berith-chain/light"
	"github.com/BerithFoundation/berith-chain/log"
	"github.com/BerithFoundation/berith-chain/node"
	"github.com/BerithFoundation/berith-chain/p2p"
	"github.com/BerithFoundation/berith-chain/p2p/discv5"
	"github.com/BerithFoundation/berith-chain/params"
	rpc "github.com/BerithFoundation/berith-chain/rpc"
)

type LightBerith struct {
	lesCommons

	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool

	// Handlers
	peers      *peerSet
	txPool     *light.TxPool
	blockchain *light.LightChain
	serverPool *serverPool
	reqDist    *requestDistributor
	retriever  *retrieveManager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *berithapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *berith.Config) (*LightBerith, error) {
	chainDb, err := berith.CreateDB(ctx, config, "lightchaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlockWithOverride(chainDb, config.Genesis, config.ConstantinopleOverride)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	stakingDB := &staking.StakingDB{}
	stakingDBPath := ctx.ResolvePath("stakingDB")
	if stkErr := stakingDB.CreateDB(stakingDBPath, staking.NewStakers); stkErr != nil {
		return nil, stkErr
	}

	lber := &LightBerith{
		lesCommons: lesCommons{
			chainDb: chainDb,
			config:  config,
			iConfig: light.DefaultClientIndexerConfig,
		},
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		peers:          peers,
		reqDist:        newRequestDistributor(peers, quitSync),
		accountManager: ctx.AccountManager,
		engine:         berith.CreateConsensusEngine(chainConfig, chainDb, stakingDB),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   berith.NewBloomIndexer(chainDb, params.BloomBitsBlocksClient, params.HelperTrieConfirmations),
	}

	lber.relay = NewLesTxRelay(peers, lber.reqDist)
	lber.serverPool = newServerPool(chainDb, quitSync, &lber.wg)
	lber.retriever = newRetrieveManager(peers, lber.reqDist, lber.serverPool)

	lber.odr = NewLesOdr(chainDb, light.DefaultClientIndexerConfig, lber.retriever)
	lber.chtIndexer = light.NewChtIndexer(chainDb, lber.odr, params.CHTFrequencyClient, params.HelperTrieConfirmations)
	lber.bloomTrieIndexer = light.NewBloomTrieIndexer(chainDb, lber.odr, params.BloomBitsBlocksClient, params.BloomTrieFrequency)
	lber.odr.SetIndexers(lber.chtIndexer, lber.bloomTrieIndexer, lber.bloomIndexer)

	// Note: NewLightChain adds the trusted checkpoint so it needs an ODR with
	// indexers already set but not started yet
	if lber.blockchain, err = light.NewLightChain(lber.odr, lber.chainConfig, lber.engine); err != nil {
		return nil, err
	}
	// Note: AddChildIndexer starts the update process for the child
	lber.bloomIndexer.AddChildIndexer(lber.bloomTrieIndexer)
	lber.chtIndexer.Start(lber.blockchain)
	lber.bloomIndexer.Start(lber.blockchain)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		lber.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	lber.txPool = light.NewTxPool(lber.chainConfig, lber.blockchain, lber.relay)
	if lber.protocolManager, err = NewProtocolManager(lber.chainConfig, light.DefaultClientIndexerConfig, true, config.NetworkId, lber.eventMux, lber.engine, lber.peers, lber.blockchain, nil, chainDb, lber.odr, lber.relay, lber.serverPool, quitSync, &lber.wg); err != nil {
		return nil, err
	}
	lber.ApiBackend = &LesApiBackend{lber, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.MinerGasPrice
	}
	lber.ApiBackend.gpo = gasprice.NewOracle(lber.ApiBackend, gpoParams)
	return lber, nil
}

func lesTopic(genesisHash common.Hash, protocolVersion uint) discv5.Topic {
	var name string
	switch protocolVersion {
	case lpv1:
		name = "LES"
	case lpv2:
		name = "LES2"
	default:
		panic(nil)
	}
	return discv5.Topic(name + "@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// Berithbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Berithbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Coinbase is the address that mining rewards will be send to (alias for Berithbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the berith package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightBerith) APIs() []rpc.API {
	return append(berithapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "berith",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "berith",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "berith",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *LightBerith) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightBerith) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightBerith) TxPool() *light.TxPool              { return s.txPool }
func (s *LightBerith) Engine() consensus.Engine           { return s.engine }
func (s *LightBerith) LesVersion() int                    { return int(ClientProtocolVersions[0]) }
func (s *LightBerith) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *LightBerith) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *LightBerith) Protocols() []p2p.Protocol {
	return s.makeProtocols(ClientProtocolVersions)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Berith protocol implementation.
func (s *LightBerith) Start(srvr *p2p.Server) error {
	log.Warn("Light client mode is an experimental feature")
	s.startBloomHandlers(params.BloomBitsBlocksClient)
	s.netRPCService = berithapi.NewPublicNetAPI(srvr, s.networkId)
	// clients are searching for the first advertised protocol in the list
	protocolVersion := AdvertiseProtocolVersions[0]
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash(), protocolVersion))
	s.protocolManager.Start(s.config.LightPeers)
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Berith protocol.
func (s *LightBerith) Stop() error {
	s.odr.Stop()
	s.bloomIndexer.Close()
	s.chtIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()
	s.engine.Close()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
