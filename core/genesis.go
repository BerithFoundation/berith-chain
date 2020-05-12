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

package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"berith-chain/berithdb"
	"berith-chain/common"
	"berith-chain/common/hexutil"
	"berith-chain/common/math"
	"berith-chain/core/rawdb"
	"berith-chain/core/state"
	"berith-chain/core/types"
	"berith-chain/log"
	"berith-chain/params"
	"berith-chain/rlp"
)

//go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Nonce      uint64              `json:"nonce"`
	Timestamp  uint64              `json:"timestamp"`
	ExtraData  []byte              `json:"extraData"`
	GasLimit   uint64              `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int            `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash         `json:"mixHash"`
	Coinbase   common.Address      `json:"coinbase"`
	Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce      math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Difficulty *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  main-net default  |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db berithdb.Database, genesis *Genesis) (*params.ChainConfig, common.Hash, error) {
	return SetupGenesisBlockWithOverride(db, genesis, nil)
}
func SetupGenesisBlockWithOverride(db berithdb.Database, genesis *Genesis, constantinopleOverride *big.Int) (*params.ChainConfig, common.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return params.MainnetChainConfig, common.Hash{}, errGenesisNoConfig
	}

	// Just commit the new block if there is no stored genesis block.
	stored := rawdb.ReadCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			log.Info("Writing default berith main-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			log.Info("Writing custom genesis block")
		}
		block, err := genesis.Commit(db)
		println("Genesis block hash", block.Header().Hash().String())
		return genesis.Config, block.Hash(), err
	}

	// Check whether the genesis block is already written.
	if genesis != nil {
		hash := genesis.ToBlock(nil).Hash()
		if hash != stored {
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
	}

	// Get the existing chain configuration.
	fmt.Println("GENESISHASH : ", stored.Hex())
	newcfg := genesis.configOrDefault(stored)
	if constantinopleOverride != nil {
		newcfg.ConstantinopleBlock = constantinopleOverride
	}
	storedcfg := rawdb.ReadChainConfig(db, stored)
	if storedcfg == nil {
		log.Warn("Found genesis block without chain config")
		rawdb.WriteChainConfig(db, stored, newcfg)
		return newcfg, stored, nil
	}
	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && stored != params.MainnetGenesisHash {
		return storedcfg, stored, nil
	}

	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := rawdb.ReadHeaderNumber(db, rawdb.ReadHeadHeaderHash(db))
	if height == nil {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	compatErr := storedcfg.CheckCompatible(newcfg, *height)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		return newcfg, stored, compatErr
	}
	rawdb.WriteChainConfig(db, stored, newcfg)
	return newcfg, stored, nil
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	case ghash == params.MainnetGenesisHash:
		return params.MainnetChainConfig
	case ghash == params.TestnetGenesisHash:
		return params.TestnetChainConfig
	default:
		return params.MainnetChainConfig
	}
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock(db berithdb.Database) *types.Block {
	if db == nil {
		db = berithdb.NewMemDatabase()
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}
	root := statedb.IntermediateRoot(false)
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
	}
	if g.GasLimit == 0 {
		head.GasLimit = params.GenesisGasLimit
	}
	if g.Difficulty == nil {
		head.Difficulty = params.GenesisDifficulty
	}
	statedb.Commit(false)
	statedb.Database().TrieDB().Commit(root, true)

	return types.NewBlock(head, nil, nil, nil)
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db berithdb.Database) (*types.Block, error) {
	block := g.ToBlock(db)
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}
	rawdb.WriteTd(db, block.Hash(), block.NumberU64(), g.Difficulty)
	rawdb.WriteBlock(db, block)
	rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), nil)
	rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
	rawdb.WriteHeadBlockHash(db, block.Hash())
	rawdb.WriteHeadHeaderHash(db, block.Hash())

	config := g.Config
	if config == nil {
		config = params.MainnetChainConfig
	}
	rawdb.WriteChainConfig(db, block.Hash(), config)
	return block, nil
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db berithdb.Database) *types.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}

// GenesisBlockForTesting creates and writes a block in which addr has the given wei balance.
func GenesisBlockForTesting(db berithdb.Database, addr common.Address, balance *big.Int) *types.Block {
	g := Genesis{Alloc: GenesisAlloc{addr: {Balance: balance}}}
	return g.MustCommit(db)
}

// DefaultGenesisBlock returns the Ethereum main net genesis block.
//[BERITH] Mainnet Genesis
func DefaultGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.MainnetChainConfig,
		Nonce:      0x00,
		Timestamp:  0x00,
		ExtraData:  hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000005559b45c940464df7145affec2f2ff4f691d92beefc1b29449332e6dd77cea91bc89db4ac9c43fa80000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   94000000,
		Difficulty: big.NewInt(1),
		Mixhash:    common.BytesToHash(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000")),
		Coinbase:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Alloc: map[common.Address]GenesisAccount{
			common.HexToAddress("Bx68dabe06325d57f2ad36cb60b754184f4f771a69"): {Balance: common.StringToBig("5000000000000000000000000000")},
		},
	}
}

// DefaultTestnetGenesisBlock returns the Ropsten network genesis block.
//[BERITH] Testnet Genesis
func DefaultTestnetGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.TestnetChainConfig,
		Nonce:      0x00,
		Timestamp:  0x00,
		ExtraData:  hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000009a17acad7ffcf6f2fc2be28ade0385d6a9d1a113b764a460e065f7bc867fa643551a18b2aad9ae3feef327f41b35a21a75dbd63eb5c1603b59494aae3c605b1c1d0d2a51d02c5fee2e0c90acbaf379b9cd10d7ec48e6730ab4778537c90b6a78a474475273966e2a0a393bcc4c5c9cc3188baea52422bae98986d57fef6100ac9fbb07438c34555260b83a85dee199749e0df598a361431ecac51efb6ad0dcbc0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:   124000000,
		Difficulty: big.NewInt(1),
		Mixhash:    common.BytesToHash(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000")),
		Coinbase:   common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Alloc: map[common.Address]GenesisAccount{
			common.HexToAddress("Bx4BEb924F14C3681CAF654eeB8684a6434BF34B73"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx942224fC3ABEF6c94DD2b45E23Bc74BDa091385c"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx21A1bF0702CCef11e86fEfa70C24c3b72f6b7a84"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxD071D21136514EbFAdd99beF2B58443Ba06F9D74"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx3ACdF7186B0FcC744D678a81C28041b3664EeA6b"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxcbf29Bd7674CA4225EA81a3bCFF4C8Cc86C669B1"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx7AA101a603c9B32bBF0a0C0Cd5d84Be7b0c165e8"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxac76025384aaF124f2D5fA3eDB00FA33d772645e"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx7c3699cF9637EEa234a9c26662eF91dE4e6727D8"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx4ae26FeFbe8EdFF6BC41B2f992aA64C8d11911bB"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxc2aa9B3376ddbE7160b38f031f5286EAaB0a5252"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx0CBaaa04408d51EeeD0e5AeBd4C0199B5831D9da"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxf282Ffa0fd84dD2fb1aC2758D028491B9231A0Fa"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx067C0F09C3805c3e664aaDF760345Ee3032f0FeC"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx8a3bd4A0089293851D5c60360560257DB644a5D2"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxd4932e16Ef9F82Aa99649Bb31862090C5546263e"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxa1E5A7EDf74bB853a14DDc5ff81DC7Bc4B5808d3"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxDda34B428944f3F0de7d57E4aEfB358E9238FA7c"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx5f49340A7F3BE8871EAc63f9649182D43D31aF1c"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx0F08191b9efA5b0d28923Cb2Ca5f12C3deD559dB"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxA20f9e5762c4f074b2e7B0BC18c013d1F4b519a0"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx84bC1C0b2F34debb7eCF90CAA0fB49253CFf2b5A"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxFB585B52833292315569115E0D30c5348E7BBb66"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx2A425A521D01B47f1B0F866d7dA6813A0E1A994f"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx36882609A5c65da99Aa7DCF1381E8820ba280812"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx3c97973646113B3148544A951a46c5Ff9Fba0177"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx757d8eEeE41A2CA7DB8Bf5c08933Ae61FE144501"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxd3D0059869Bab7A8B284dF853009aadBa62e2af5"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx3B1Af99F36a41D1Edc863426C77fcf43bfFDd662"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxA2c19Eb00C6219A21ecdbd044E34505ECAdf8513"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx3c3fe9EC3F10334a3129162EBCc21503E4f5827E"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx27d100Ecc1BDfE33c69405b0B42351a6f364B456"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx3A360E5dc66f022AD5E2c66bd5eaF19ce5527061"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxB72817885Cdb92967BC5958717178084BA67fd39"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx6fdd2138806762daA9118E2289f5A2690206f9D0"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx8CdE6c838f2dD49266Ece51437f621E9C2796d6a"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx000f58Bfa9D30eB8c330920d91f1AF27c3116184"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx56af00E820408e56dc5cC09f13A921E35f78665A"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxd067f556eE3Ed5053e259c78439dDd746d3617B5"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxfE83127cCe0D7B6434a3379B255370eCa9e9360B"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx3132c45acFA963D22c26AeE7859B251f74fC2Acb"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx7E80b2AFb8DEdd8Af6B3B8D08a04B32eaB1f87dB"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx099f03Df3684570f78E971Afe65Eb6Ab6ee927A4"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxE746c98d22BB9A01F9B203e91e58f0dfC44F2FDc"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxA99719496afe5cD7D9C2dF7dd9dd615A6Fe2eD61"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx45af910D5253326088902663F5c5F1Bb8342dd36"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxd1d06E914642D4ED08924F467eB35f9496801015"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx79c752dCac1247a818F446E4C31edF83b0Bd71BE"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx461B5A71E4cCDDe9A36940fE6d823b346f2Ca273"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx2e501AD9bc16ED20578907eFd3012322C88c15Dc"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx8f23fF1Eb99cC31dF21C05DDDA0c79E27e9AF5c6"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx38BdDe0921Fb3bE1586Ae203b3f9E3894e48f385"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx86a09B4099e7244F330A546D51845EC0015F6697"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxefa4BB3590e0c3415Cf1d4AA826496a2c755b07c"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxBc25C9167E61f04e4B4a4d623B9A8BdA1E1C6D32"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxf7c99180d9A303aDb65D8b7B5aFacb269860d14e"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx94367098a583E5bE4bc2E88E9aD1DeD6EFfDC956"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx1dfe3ffCA507111875f9635984bD5D1858c2E424"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx7A7D0F64852Acbf49bfABF7525457b55b44676D3"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxD3d4553f3a50d645e3A5c8553cE1785be32911E7"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx4dD7bF198A7a2525A240E58E0B0f703688816689"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx1c1A044F72a7d866BdfC241B72ad5E5444De5F8F"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxf88FC73eE8A7E29fF21F9C0e316FE524241c95A7"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxA3672E627d42ae2071e8D899800adFB401047326"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxB3b5ade64f423820E5C631813d07d5285Ebad001"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxEdB6DcefC8F81EC66C50C7cBaf2C0965359cb168"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx084482977AeB3A884c3cA6fD66FcA031437ac386"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxdB12288F4808D588c5Ac980603186F83e2870f7f"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx79bBE43ae8198CaEb47A0728aBc25a7741dd9A83"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxBD1b89ce210B493c2786736f39a66f5c5b4E76f8"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx44d7d228284E65aD9bF006b3f03B40E3C41bb120"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxd7fdb9cE93F40F93F3D9a260A5F6b1D96093A81C"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxCa154af5b078B0Dde927Ae7bFa88eBb3C1273BC3"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxCE1B04C45513252CE262b0a16a2F915ea47F6234"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxfb41e47a1Fabe06005d6B894519a0FBA306aFb3A"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx07c11f69A5341b1bD215b1ecF682a0b2aAB3da1B"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx2d38663a25D5672118f12b9567d3cd49941A0863"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx8f76f014Ab2BCcBfBEb030e3Ae990330f566F296"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx5D73700d5841f8d4B7Be6A72D088FBd8cf15E90E"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx134733df5E4dd18BFeA1dAB264AABa8470bDD8D4"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx40DE707851676311A5e47e560738980a9Ee7Abfe"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx523F7e4A6eE59F3223Aa0B94250c1C0043230D4d"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx5a9938cD3EF1bcC77c4dB22781E75eB8Ed5002d7"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx62aEdd26F6D47fF60BcCE095D974d51A7c0DE607"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxb3d88410EB35543F8561a773CA3E6556994E915d"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx318A83d9D879A7C4bF920CB04684F1Ad2e7a523F"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx3cc746Cd15d0c7eE5d57c3E72BCb6Be58fB40F40"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx12fA7791a957ba8EAabbC884A7eE1e5094a731eb"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx92C545769B4CF82f4426a090f7fB27a921f2330f"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxbB5D486FD9d585d35e866357C6761dc187715173"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx2BcB3d0DA694Fcc252cb7813Ad8D0D189Df74960"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx4E16f406c1011d24E31d5A27a7BF1F0A7f993d77"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxB35499Bd1D40488cb844D61276805593749a6497"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx87d07e3C3F06ba985804ddACc0B263805d1ffa30"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx081c5bAAF5a4d7B38EA2c1C18a1D288a14e1F4BA"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxEd31C92c9b226fB08b651E5CF7bF8b4BCB870121"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx5058C89505c7b49DD94EF36C059Dd112773C58eA"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx0Ddb26d87E077Dfa3591BA69d77794374144b7B5"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxa8966eaBC3dA07Cc34B3A68caA92224638065184"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxd693373f0aFc8E2520C0543785903c36F3213b4E"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx30b68A16Ea00F56f0099DE97145cDfd157425025"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx6686951d0D9d0bcAa17C590a6F918AF444B850C4"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxda151361eEa66f290cE1D800fEeBBD48C61ae60c"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx5C73D49e277EC8AaB1baB8aCd68761DFcCa433b9"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("BxdC51ECc276f496c0be08A8fa2B4521CEac8A3491"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bxb99bA521aD7AFB453549180E04FcC9f10602fe55"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx1D89364e05c84cF570Ba62Db1f645c9145887971"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx1afD46a6118d47a43fD793D364b18fE8639D8eE6"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx5696542Bf9a4d0e5a4cCA2f2F306c23A294197dd"): {Balance: common.StringToBig("100000000000000000000000000")},
			common.HexToAddress("Bx9a17acad7ffcf6f2fc2be28ade0385d6a9d1a113"): {Balance: common.StringToBig("10000000000000000000000000000")},
			common.HexToAddress("Bxb764a460e065f7bc867fa643551a18b2aad9ae3f"): {Balance: common.StringToBig("10000000000000000000000000000")},
			common.HexToAddress("Bxeef327f41b35a21a75dbd63eb5c1603b59494aae"): {Balance: common.StringToBig("10000000000000000000000000000")},
			common.HexToAddress("Bx3c605b1c1d0d2a51d02c5fee2e0c90acbaf379b9"): {Balance: common.StringToBig("10000000000000000000000000000")},
			common.HexToAddress("Bxcd10d7ec48e6730ab4778537c90b6a78a4744752"): {Balance: common.StringToBig("10000000000000000000000000000")},
			common.HexToAddress("Bx73966e2a0a393bcc4c5c9cc3188baea52422bae9"): {Balance: common.StringToBig("10000000000000000000000000000")},
			common.HexToAddress("Bx8986d57fef6100ac9fbb07438c34555260b83a85"): {Balance: common.StringToBig("10000000000000000000000000000")},
			common.HexToAddress("Bxdee199749e0df598a361431ecac51efb6ad0dcbc"): {Balance: common.StringToBig("10000000000000000000000000000")},
		},
	}
}

func decodePrealloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}
