// Copyright 2014 The go-ethereum Authors
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

// Package berith implements the Berith protocol.
package berith

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/BerithFoundation/berith-chain/berith/staking"

	"github.com/BerithFoundation/berith-chain/consensus/bsrr"

	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/berith/brtapi"
	"github.com/BerithFoundation/berith-chain/berith/downloader"
	"github.com/BerithFoundation/berith-chain/berith/filters"
	"github.com/BerithFoundation/berith-chain/berith/gasprice"
	"github.com/BerithFoundation/berith-chain/berith/stakingdb"
	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/common/hexutil"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/core"
	"github.com/BerithFoundation/berith-chain/core/bloombits"
	"github.com/BerithFoundation/berith-chain/core/rawdb"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/core/vm"
	"github.com/BerithFoundation/berith-chain/event"
	"github.com/BerithFoundation/berith-chain/internal/berithapi"
	"github.com/BerithFoundation/berith-chain/log"
	"github.com/BerithFoundation/berith-chain/miner"
	"github.com/BerithFoundation/berith-chain/node"
	"github.com/BerithFoundation/berith-chain/p2p"
	"github.com/BerithFoundation/berith-chain/params"
	"github.com/BerithFoundation/berith-chain/rlp"
	"github.com/BerithFoundation/berith-chain/rpc"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

/*
[BERITH]
Berith implements the Berith full node service.
*/
type Berith struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the Berith

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDb berithdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	APIBackend *BerAPIBackend

	miner      *miner.Miner
	gasPrice   *big.Int
	berithbase common.Address

	networkID     uint64
	netRPCService *berithapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and berithbase)

	/*
	[BERITH]
	database for staker infos
	*/
	stakingDB *stakingdb.StakingDB
}

func (s *Berith) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new Berith object (including the
// initialisation of the common Berith object)
func New(ctx *node.ServiceContext, config *Config) (*Berith, error) {
	// Ensure configuration values are compatible and sane
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run berith.Berith in light sync mode, use les.LightBerith")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	if config.MinerGasPrice == nil || config.MinerGasPrice.Cmp(common.Big0) <= 0 {
		log.Warn("Sanitizing invalid miner gas price", "provided", config.MinerGasPrice, "updated", DefaultConfig.MinerGasPrice)
		config.MinerGasPrice = new(big.Int).Set(DefaultConfig.MinerGasPrice)
	}
	// Assemble the Berith object
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlockWithOverride(chainDb, config.Genesis, config.ConstantinopleOverride)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	stakingDB := &stakingdb.StakingDB{}
	stakingDBPath := ctx.ResolvePath("stakingDB")
	if stkErr := stakingDB.CreateDB(stakingDBPath, staking.NewStakers); stkErr != nil {
		return nil, stkErr
	}
	engine := CreateConsensusEngine(chainConfig, chainDb, stakingDB)
	ber := &Berith{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         engine,
		shutdownChan:   make(chan bool),
		networkID:      config.NetworkId,
		gasPrice:       config.MinerGasPrice,
		berithbase:     config.Berithbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks, params.BloomConfirms),
		stakingDB:      stakingDB,
	}

	log.Info("Initialising Berith protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d).\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}
	var (
		vmConfig = vm.Config{
			EnablePreimageRecording: config.EnablePreimageRecording,
			EWASMInterpreter:        config.EWASMInterpreter,
			EVMInterpreter:          config.EVMInterpreter,
		}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieCleanLimit: config.TrieCleanCache, TrieDirtyLimit: config.TrieDirtyCache, TrieTimeLimit: config.TrieTimeout}
	)
	ber.blockchain, err = core.NewBlockChain(stakingDB, chainDb, cacheConfig, ber.chainConfig, ber.engine, vmConfig, ber.shouldPreserve)
	if err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		ber.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	ber.bloomIndexer.Start(ber.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	ber.txPool = core.NewTxPool(config.TxPool, ber.chainConfig, ber.blockchain)

	if ber.protocolManager, err = NewProtocolManager(ber.chainConfig, config.SyncMode, config.NetworkId, ber.eventMux, ber.txPool, ber.engine, ber.blockchain, chainDb, config.Whitelist); err != nil {
		return nil, err
	}

	ber.miner = miner.New(ber, ber.chainConfig, ber.EventMux(), ber.engine, config.MinerRecommit, config.MinerGasFloor, config.MinerGasCeil, ber.isLocalBlock)
	ber.miner.SetExtra(makeExtraData(config.MinerExtraData))

	ber.APIBackend = &BerAPIBackend{ber, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.MinerGasPrice
	}
	ber.APIBackend.gpo = gasprice.NewOracle(ber.APIBackend, gpoParams)

	return ber, nil
}

func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"berith",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (berithdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*berithdb.LDBDatabase); ok {
		db.Meter("berith/db/chaindata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an Berith service
func CreateConsensusEngine(chainConfig *params.ChainConfig, db berithdb.Database, stakingDB *stakingdb.StakingDB) consensus.Engine {
	return bsrr.NewCliqueWithStakingDB(stakingDB, chainConfig.Bsrr, db)
}

// APIs return the collection of RPC services the berith package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
/*
[BERITH]
Registration of each service implementation
services : miner, admin, berith, etc...
*/
func (s *Berith) APIs() []rpc.API {
	apis := berithapi.GetAPIs(s.APIBackend)
	apis = append(apis, brtapi.GetAPIs(s.APIBackend, s.miner)...)
	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "berith",
			Version:   "1.0",
			Service:   NewPublicBerithAPI(s),
			Public:    true,
		}, {
			Namespace: "berith",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "berith",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "berith",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *Berith) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Berith) Berithbase() (eb common.Address, err error) {
	s.lock.RLock()
	berithbase := s.berithbase
	s.lock.RUnlock()

	if berithbase != (common.Address{}) {
		return berithbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			berithbase := accounts[0].Address

			s.lock.Lock()
			s.berithbase = berithbase
			s.lock.Unlock()

			log.Info("Berithbase automatically configured", "address", berithbase)
			return berithbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("berithbase must be explicitly specified")
}

// isLocalBlock checks whether the specified block is mined
// by local miner accounts.
//
// We regard two types of accounts as local miner account: berithbase
// and accounts specified via `txpool.locals` flag.
func (s *Berith) isLocalBlock(block *types.Block) bool {
	author, err := s.engine.Author(block.Header())
	if err != nil {
		log.Warn("Failed to retrieve block author", "number", block.NumberU64(), "hash", block.Hash(), "err", err)
		return false
	}
	// Check whether the given address is berithbase.
	s.lock.RLock()
	berithbase := s.berithbase
	s.lock.RUnlock()
	if author == berithbase {
		return true
	}
	// Check whether the given address is specified by `txpool.local`
	// CLI flag.
	for _, account := range s.config.TxPool.Locals {
		if account == author {
			return true
		}
	}
	return false
}

// shouldPreserve checks whether we should preserve the given block
// during the chain reorg depending on whether the author of block
// is a local account.
func (s *Berith) shouldPreserve(block *types.Block) bool {
	return s.isLocalBlock(block)
}

// SetBerithbase sets the mining reward address.
func (s *Berith) SetBerithbase(base common.Address) {
	s.lock.Lock()
	s.berithbase = base
	s.lock.Unlock()

	s.miner.SetBerithbase(base)
}

// StartMining starts the miner with the given number of CPU threads. If mining
// is already running, this method adjust the number of threads allowed to use
// and updates the minimum price required by the transaction pool.
func (s *Berith) StartMining(threads int) error {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		log.Info("Updated mining threads", "threads", threads)
		if threads == 0 {
			threads = -1 // Disable the miner from within
		}
		th.SetThreads(threads)
	}
	// If the miner was not running, initialize it
	if !s.IsMining() {
		// Propagate the initial price point to the transaction pool
		s.lock.RLock()
		price := s.gasPrice
		s.lock.RUnlock()
		s.txPool.SetGasPrice(price)

		// Configure the local mining address
		eb, err := s.Berithbase()
		if err != nil {
			log.Error("Cannot start mining without berithbase", "err", err)
			return fmt.Errorf("berithbase missing: %v", err)
		}
		if bsrr, ok := s.engine.(*bsrr.BSRR); ok {
			wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
			if wallet == nil || err != nil {
				log.Error("Berithbase account unavailable locally", "err", err)
				return fmt.Errorf("signer missing: %v", err)
			}
			bsrr.Authorize(eb, wallet.SignHash)
		}
		// If mining is started, we can disable the transaction rejection mechanism
		// introduced to speed sync times.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)

		go s.miner.Start(eb)
	}
	return nil
}

// StopMining terminates the miner, both at the consensus engine level as well as
// at the block creation level.
func (s *Berith) StopMining() {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	// Stop the block creating itself
	s.miner.Stop()
}

func (s *Berith) IsMining() bool      { return s.miner.Mining() }
func (s *Berith) Miner() *miner.Miner { return s.miner }

func (s *Berith) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Berith) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Berith) TxPool() *core.TxPool               { return s.txPool }
func (s *Berith) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Berith) Engine() consensus.Engine           { return s.engine }
func (s *Berith) ChainDb() berithdb.Database         { return s.chainDb }
func (s *Berith) IsListening() bool                  { return true } // Always listening
func (s *Berith) BerVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Berith) NetVersion() uint64                 { return s.networkID }
func (s *Berith) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *Berith) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Berith protocol implementation.
func (s *Berith) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers(params.BloomBitsBlocks)

	// Start the RPC service
	s.netRPCService = berithapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= srvr.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, srvr.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Berith protocol.
func (s *Berith) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.engine.Close()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)
	return nil
}
