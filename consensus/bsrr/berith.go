/*
d8888b. d88888b d8888b. d888888b d888888b db   db
88  `8D 88'     88  `8D   `88'   `~~88~~' 88   88
88oooY' 88ooooo 88oobY'    88       88    88ooo88
88~~~b. 88~~~~~ 88`8b      88       88    88~~~88
88   8D 88.     88 `88.   .88.      88    88   88
Y8888P' Y88888P 88   YD Y888888P    YP    YP   YP

	  copyrights by ibizsoftware 2018 - 2019
*/

/**
[BERITH]
- 합의 알고리즘 인터페이스 구현체로 Berith 합의 절차를 여기서 처리함
- 해더 검증및 바디 데이터 검증을 함
- 바디 데이터 검증
  BC 체크, 그룹체크, 우선순위 검증
- 바디의 Tx를 확인 하여 Staking DB 에 기록 하고 선출
**/

package bsrr

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/BerithFoundation/berith-chain/rpc"

	"github.com/BerithFoundation/berith-chain/accounts"
	"github.com/BerithFoundation/berith-chain/berith/selection"
	"github.com/BerithFoundation/berith-chain/berith/staking"
	"github.com/BerithFoundation/berith-chain/berithdb"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/consensus"
	"github.com/BerithFoundation/berith-chain/consensus/misc"
	"github.com/BerithFoundation/berith-chain/core/state"
	"github.com/BerithFoundation/berith-chain/core/types"
	"github.com/BerithFoundation/berith-chain/crypto"
	"github.com/BerithFoundation/berith-chain/crypto/sha3"
	"github.com/BerithFoundation/berith-chain/log"
	"github.com/BerithFoundation/berith-chain/params"
	"github.com/BerithFoundation/berith-chain/rlp"
	lru "github.com/hashicorp/golang-lru"
)

const (
	inmemorySnapshots  = 128     // Number of recent vote snapshots to keep in memory
	inmemorySigners    = 128 * 3 // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096    // Number of recent block signatures to keep in memory

	termDelay  = 100 * time.Millisecond // Delay per signer in the same group
	groupDelay = 1 * time.Second        // Delay per groups

	commonDiff = 3 // A constant that specifies the maximum number of people in a group when dividing a signer's candidates into multiple groups
)

var (
	RewardBlock  = big.NewInt(500)
	StakeMinimum = new(big.Int).Mul(big.NewInt(100000), common.UnitForBer)
	SlashRound   = uint64(2)
	ForkFactor   = 1.0

	epochLength = uint64(360) // Default number of blocks after which to checkpoint and reset the pending votes

	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.

	diffWithoutStaker = int64(1234)
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	// errExtraSigners is returned if non-checkpoint block contain signer data in
	// their extra-data fields.
	errExtraSigners = errors.New("non-checkpoint block contains extra signer list")

	// errInvalidCheckpointSigners is returned if a checkpoint block contains an
	// invalid list of signers (i.e. non divisible by 20 bytes).
	errInvalidCheckpointSigners = errors.New("invalid signer list on checkpoint block")

	// errInvalidMixDigest is returned if a block's mix digest is non-zero.
	errInvalidMixDigest = errors.New("non-zero mix digest")

	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")

	// errInvalidDifficulty is returned if the difficulty of a block neither 1 or 2.
	errInvalidDifficulty = errors.New("invalid difficulty")

	// ErrInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")

	// errUnauthorizedSigner is returned if a header is signed by a non-authorized entity.
	errUnauthorizedSigner = errors.New("unauthorized signer")

	errNoData = errors.New("no data")

	// errInvalidNonce is returned if a nonce is less than or equals to 0.
	errInvalidNonce = errors.New("invalid nonce")

	errStakingList = errors.New("not found staking list")

	errMissingState = errors.New("state missing")

	errBIP1 = errors.New("error when fork network to BIP1")
)

// SignerFn is a signer callback function to request a hash to be signed by a
// backing account.
type SignerFn func(accounts.Account, []byte) ([]byte, error)

// sigHash returns the hash which is used as input for the proof-of-authority
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	_ = rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])
	return hash
}

// ecrecover extracts the Berith account address from a signed header.
func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	// If the signature's already cached, return that
	hash := header.Hash()
	if address, known := sigcache.Get(hash); known {
		return address.(common.Address), nil
	}
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Berith address
	pubkey, err := crypto.Ecrecover(sigHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, signer)
	return signer, nil
}

type BSRR struct {
	config *params.BSRRConfig // Consensus engine configuration parameters
	db     berithdb.Database  // Database to store and retrieve snapshot checkpoints
	//[BERITH] stakingDB clique 구조체에 추가
	stakingDB staking.DataBase //stakingList를 저장하는 DB
	cache     *lru.ARCCache    //stakingList를 저장하는 cache

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	signer common.Address // Berith address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with
	lock   sync.RWMutex   // Protects the signer fields

	proposals map[common.Address]bool // Current list of proposals we are pushing

	// The fields below are for testing only
	rankGroup common.SequenceGroup // grouped by rank
}

//[BERITH]
//New 새로운 BSRR 구조체를 만드는 함수
func New(config *params.BSRRConfig, db berithdb.Database) *BSRR {
	conf := config
	if conf.Epoch == 0 {
		conf.Epoch = epochLength
	}

	if conf.Rewards != nil {
		if conf.Rewards.Cmp(big.NewInt(0)) == 0 {
			conf.Rewards = RewardBlock
		}
	} else {
		conf.Rewards = RewardBlock
	}

	if conf.StakeMinimum != nil {
		if conf.StakeMinimum.Cmp(big.NewInt(0)) == 0 {
			conf.StakeMinimum = StakeMinimum
		}
	} else {
		conf.StakeMinimum = StakeMinimum
	}

	if conf.SlashRound != 0 {
		if conf.SlashRound == 0 {
			conf.SlashRound = 1
		}
	} else {
		conf.SlashRound = SlashRound
	}

	if conf.ForkFactor <= 0.0 || conf.ForkFactor > 1.0 {
		conf.ForkFactor = ForkFactor
	}

	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)
	//[BERITH] 캐쉬 인스턴스 생성및 사이즈 지정
	cache, _ := lru.NewARC(inmemorySigners)

	return &BSRR{
		config:     conf,
		db:         db,
		recents:    recents,
		signatures: signatures,
		cache:      cache,
		proposals:  make(map[common.Address]bool),
		rankGroup:  &common.ArithmeticGroup{CommonDiff: commonDiff},
	}
}

//[BERITH]
//NewCliqueWithStakingDB StakingDB를 받아 새로운 BSRR 구조체를 생성하는 함수
func NewCliqueWithStakingDB(stakingDB staking.DataBase, config *params.BSRRConfig, db berithdb.Database) *BSRR {
	engine := New(config, db)
	engine.stakingDB = stakingDB
	// Synchronize the engine.config and chainConfig.
	return engine
}

// Author implements consensus.Engine, returning the Berith address recovered
// from the signature in the header's extra-data section.
func (c *BSRR) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header, c.signatures)
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (c *BSRR) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return c.verifyHeader(chain, header, nil)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (c *BSRR) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := c.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (c *BSRR) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()

	// Don't waste time checking blocks from the future
	if header.Time.Cmp(big.NewInt(time.Now().Unix())) > 0 {
		return consensus.ErrFutureBlock
	}
	// Checkpoint blocks need to enforce zero beneficiary
	checkpoint := (number % c.config.Epoch) == 0

	// Check that the extra-data contains both the vanity and signature
	if len(header.Extra) < extraVanity {
		return errMissingVanity
	}
	if len(header.Extra) < extraVanity+extraSeal {
		return errMissingSignature
	}
	// Ensure that the extra-data contains a signer list on checkpoint, but none otherwise
	signersBytes := len(header.Extra) - extraVanity - extraSeal
	if !checkpoint && signersBytes != 0 {
		return errExtraSigners
	}
	if checkpoint && signersBytes%common.AddressLength != 0 {
		return errInvalidCheckpointSigners
	}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (common.Hash{}) {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if header.UncleHash != uncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block nonce is greater than 0
	if number > 0 && header.Nonce.Uint64() < 1 {
		return errInvalidNonce
	}

	// If all checks passed, validate any special fields for hard forks
	if err := misc.VerifyForkHashes(chain.Config(), header, false); err != nil {
		return err
	}
	// All basic checks passed, verify cascading fields
	return c.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (c *BSRR) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time.Uint64()+c.config.Period > header.Time.Uint64() {
		return ErrInvalidTimestamp
	}

	// All basic checks passed, verify the seal and return
	return c.verifySeal(chain, header, parents)
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (c *BSRR) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (c *BSRR) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return c.verifySeal(chain, header, nil)
}

// verifySeal checks whether the signature contained in the header satisfies the
// consensus protocol requirements. The method accepts an optional list of parent
// headers that aren't yet part of the local blockchain to generate the snapshots
// from.
/*
	[Berith]
	verifySeal method is necessary to implement Engine interface but not used.
	The logic that verifies the signature contained in the header is in the Finalize method.
*/
func (c *BSRR) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number
	if number.Uint64() == 0 {
		return errUnknownBlock
	}

	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (c *BSRR) Prepare(chain consensus.ChainReader, header *types.Header) error {
	header.Nonce = types.BlockNonce{}
	number := header.Number.Uint64()

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	target, exist := c.getStakeTargetBlock(chain, parent)
	if !exist {
		return consensus.ErrUnknownAncestor
	}

	// Set the correct difficulty and nonce
	diff, rank := c.calcDifficultyAndRank(c.signer, chain, 0, target)
	if rank < 1 {
		return errUnauthorizedSigner
	}
	header.Difficulty = diff
	// nonce is used to check order of staking list
	header.Nonce = types.EncodeNonce(uint64(rank))

	// Ensure the extra data has all it's components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]

	header.Extra = append(header.Extra, make([]byte, extraSeal)...)

	// Mix digest is reserved for now, set to empty
	header.MixDigest = common.Hash{}

	header.Time = new(big.Int).Add(parent.Time, new(big.Int).SetUint64(c.config.Period))
	if header.Time.Int64() < time.Now().Unix() {
		header.Time = big.NewInt(time.Now().Unix())
	}
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (c *BSRR) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	//[Berith] 부모블록의 StakingList를 얻어온다.
	var stks staking.Stakers
	stks, err := c.getStakers(chain, header.Number.Uint64()-1, header.ParentHash)
	if err != nil {
		return nil, errStakingList
	}

	if header.Coinbase != common.HexToAddress("0") {
		var signers signers

		parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
		if parent == nil {
			log.Warn("unknown ancestor", "parent", "nil")
		}

		if chain.Config().IsBIP1Block(header.Number) {
			stks, err = c.supportBIP1(chain, parent, stks)
			if err != nil {
				return nil, errBIP1
			}
		}
		target, exist := c.getStakeTargetBlock(chain, parent)
		if !exist {
			return nil, consensus.ErrUnknownAncestor
		}

		signers, err := c.getSigners(chain, target)
		if err != nil {
			return nil, errUnauthorizedSigner
		}

		signerMap := signers.signersMap()
		if _, ok := signerMap[header.Coinbase]; !ok {
			return nil, errUnauthorizedSigner
		}

		predicted, rank := c.calcDifficultyAndRank(header.Coinbase, chain, 0, target)
		if rank < 1 {
			return nil, errUnauthorizedSigner
		}

		if predicted.Cmp(header.Difficulty) != 0 {
			return nil, errInvalidDifficulty
		}
		if header.Nonce.Uint64() != uint64(rank) {
			return nil, errInvalidNonce
		}
	}

	//[BERITH] 전달받은 블록의 트랜잭션을 정보를 토대로 StateDB의 데이터를 수정한다.
	err = c.setStakersWithTxs(state, chain, stks, txs, header)
	if err != nil {
		return nil, errStakingList
	}

	//Reward 보상
	c.accumulateRewards(chain, state, header)

	//[BERITH] 수정된 StateDB의 데이터를 commit한다.
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (c *BSRR) Authorize(signer common.Address, signFn SignerFn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.signer = signer
	c.signFn = signFn
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (c *BSRR) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	header := block.Header()

	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// Don't hold the signer fields for the entire sealing procedure
	c.lock.RLock()
	signer, signFn := c.signer, c.signFn
	c.lock.RUnlock()

	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// Checks target block and signers
	target, exist := c.getStakeTargetBlock(chain, parent)
	if !exist {
		return consensus.ErrUnknownAncestor
	}

	signers, err := c.getSigners(chain, target)
	if err != nil {
		return err
	}
	if _, authorized := signers.signersMap()[signer]; !authorized {
		return errUnauthorizedSigner
	}

	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := time.Unix(header.Time.Int64(), 0).Sub(time.Now()) // nolint: gosimple
	_, rank := c.calcDifficultyAndRank(header.Coinbase, chain, 0, target)
	if rank == -1 {
		return errUnauthorizedSigner
	}

	//delay += c.getDelay(rank)
	temp, err := c.getDelay(rank)
	if err != nil {
		return err
	}
	delay += temp

	// Sign all the things!
	sighash, err := signFn(accounts.Account{Address: signer}, sigHash(header).Bytes())
	if err != nil {
		return err
	}
	copy(header.Extra[len(header.Extra)-extraSeal:], sighash)
	// Wait until sealing is terminated or delay timeout.
	log.Trace("Waiting for slot to sign and propagate", "delay", common.PrettyDuration(delay))
	go func() {
		select {
		case <-stop:
			return
		case <-time.After(delay):
		}

		select {
		case results <- block.WithSeal(header):
		default:
			log.Warn("Sealing result is not read by miner", "sealhash", c.SealHash(header))
		}
	}()
	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (c *BSRR) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	//target, exist := c.getAncestor(chain, int64(c.config.Epoch), parent)
	target, exist := c.getStakeTargetBlock(chain, parent)
	if !exist {
		return big.NewInt(0)
	}
	diff, _ := c.calcDifficultyAndRank(c.signer, chain, time, target)
	return diff
}

func (c *BSRR) getAncestor(chain consensus.ChainReader, n int64, header *types.Header) (*types.Header, bool) {
	target := header
	targetNumber := new(big.Int).Sub(header.Number, big.NewInt(n))
	for target != nil && target.Number.Cmp(big.NewInt(0)) > 0 && target.Number.Cmp(targetNumber) > 0 {
		target = chain.GetHeader(target.ParentHash, target.Number.Uint64()-1)
	}

	if target == nil {
		return &types.Header{}, false
	}

	return target, chain.HasBlockAndState(target.Hash(), target.Number.Uint64())
}

// [BERITH] getStakeTargetBlock 주어진 parent header에 대하여 miner를 결정 할 target block을 반환한다.
// 1) [0 ~ epoch-1]     : target == 블록 넘버 0(즉, genesis block) 인 블록
// 2) [epoch ~ 2epoch-1] : target == 블록 넘버 epoch 인 블록
// 3) [2epoch ~ ...)       : target == 블록 넘버 - epoch 인 블록
func (c *BSRR) getStakeTargetBlock(chain consensus.ChainReader, parent *types.Header) (*types.Header, bool) {
	if parent == nil {
		return &types.Header{}, false
	}

	var targetNumber uint64
	blockNumber := parent.Number.Uint64()
	d := blockNumber / c.config.Epoch

	if d > 1 {
		return c.getAncestor(chain, int64(c.config.Epoch), parent)
	}

	switch d {
	case 0:
		targetNumber = 0
	case 1:
		targetNumber = c.config.Epoch
	}

	target := chain.GetHeaderByNumber(targetNumber)
	if target != nil {
		return target, chain.HasBlockAndState(target.Hash(), targetNumber)
	}
	return target, false
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *BSRR) SealHash(header *types.Header) common.Hash {
	return sigHash(header)
}

// [BERITH] calcDifficultyAndRank 주어진 address에 대하여 블록을 생성할 때의 난이도, 순위를 반환하는 메서드
// 1) [0, epoch] -> genesis block의 extra data에서 추출 후
// ==> (1234,1) or (0, -1) 반환
// 2) [epoch+1, ~) -> target의 블록까지 존재하는 스테이킹 리스트기반 (diff, rank) 반환
// ==> (diff,rank) or (0, -1) 반환
func (c *BSRR) calcDifficultyAndRank(signer common.Address, chain consensus.ChainReader, time uint64, target *types.Header) (*big.Int, int) {
	// extract diff and rank from genesis's extra data
	if target.Number.Cmp(big.NewInt(0)) == 0 {
		log.Info("default difficulty and rank", "diff", diffWithoutStaker, "rank", 1)
		return big.NewInt(diffWithoutStaker), 1
	}

	stks, err := c.getStakers(chain, target.Number.Uint64(), target.Hash())

	if err != nil {
		log.Error("failed to get stakers", "err", err.Error())
		return big.NewInt(0), -1
	}

	stateDB, err := chain.StateAt(target.Root)

	if err != nil {
		log.Error("failed to get state", "err", err.Error())
		return big.NewInt(0), -1
	}

	results := selection.SelectBlockCreator(chain.Config(), target.Number.Uint64(), target.Hash(), stks, stateDB)

	max := c.getMaxMiningCandidates(len(results))

	if results[signer].Rank > max {
		log.Warn("out of rank", "hash", target.Hash().Hex(), "rank", results[signer].Rank, "max", max)
		return big.NewInt(0), -1
	}

	return results[signer].Score, results[signer].Rank
}

// getDelay 주어진 rank에 따라 블록 Sealing에 대한 지연 시간을 반환한다.
// 항상 0보다 크거나 같은 값을 반환
func (c *BSRR) getDelay(rank int) (time.Duration, error) {
	if rank <= 1 {
		return time.Duration(0), nil
	}

	// 각 그룹별 지연시간
	groupOrder, err := c.rankGroup.GetGroupOrder(rank)
	if err != nil {
		return time.Duration(0), err
	}
	delay := time.Duration(groupOrder-1) * groupDelay

	// 그룹 내 지연 시간
	startRank, _, err := c.rankGroup.GetGroupRange(groupOrder)
	if err != nil {
		return time.Duration(0), err
	}
	delay += time.Duration(rank-startRank) * termDelay

	return delay, nil
}

// Close implements consensus.Engine. It's a noop for clique as there are no background threads.
func (c *BSRR) Close() error {
	return nil
}

func getReward(config *params.ChainConfig, header *types.Header) *big.Int {
	const (
		defaultBlockCreationSec    = 10      // Blocks are created every 10 seconds by default.
		blockNumberAt1Year         = 3150000 // If a block is created every 10 seconds, this number of the block created at the time of 1 year.
		defaultReward              = 26      // The basic reward is 26 tokens.
		additionalReward           = 5       // Additional rewards are paid for one year.
		blockSectionDivisionNumber = 7370000 // Reference value for dividing a block into 50 sections
		groupingValue              = 0.5     // Constant for grouping two groups to have the same Reward Subtract
	)

	number := header.Number.Uint64()
	// Reward after a specific block
	if number < config.Bsrr.Rewards.Uint64() {
		return big.NewInt(0)
	}

	// Value to correct Reward when block creation time is changed.
	correctionValue := float64(config.Bsrr.Period) / defaultBlockCreationSec
	correctedBlockNumber := float64(number) * correctionValue

	var addtional float64 = 0
	if correctedBlockNumber <= blockNumberAt1Year {
		addtional = additionalReward
	}

	/*
		[Berith]
		The reward payment decreases as the time increases, and for this purpose, the block is divided into 50 sections.
		The same amount is deducted for every two sections.
	*/
	reward := (defaultReward - math.Round(correctedBlockNumber/blockSectionDivisionNumber)*groupingValue + addtional) * correctionValue
	if reward <= 0 {
		return big.NewInt(0)
	}
	temp := reward * 1e+10
	return new(big.Int).Mul(big.NewInt(int64(temp)), big.NewInt(1e+8))
}

// AccumulateRewards credits the coinbase of the given block with the mining
// reward.
func (c *BSRR) accumulateRewards(chain consensus.ChainReader, state *state.StateDB, header *types.Header) {
	config := chain.Config()
	state.AddBehindBalance(header.Coinbase, header.Number, getReward(config, header))

	//과거 시점의 블록 생성자 가져온다.
	target, exist := c.getAncestor(chain, int64(config.Bsrr.Epoch), header)
	if !exist {
		return
	}

	signers, err := c.getSigners(chain, target)
	if err != nil {
		return
	}

	//all node block result
	for _, addr := range signers {
		behind, err := state.GetFirstBehindBalance(addr)
		if err != nil {
			continue
		}

		target := new(big.Int).Add(behind.Number, new(big.Int).SetUint64(config.Bsrr.Epoch))
		if header.Number.Cmp(target) == -1 {
			continue
		}

		if behind.Balance.Cmp(new(big.Int).SetInt64(int64(0))) != 1 {
			continue
		}

		//bihind --> main
		state.AddBalance(addr, behind.Balance)

		state.RemoveFirstBehindBalance(addr)
	}
}

func (c *BSRR) supportBIP1(chain consensus.ChainReader, parent *types.Header, stks staking.Stakers) (staking.Stakers, error) {
	st, err := chain.StateAt(parent.Root)
	if err != nil {
		return nil, err
	}

	for _, addr := range stks.AsList() {
		if st.GetStakeBalance(addr).Cmp(c.config.StakeMinimum) < 0 {
			stks.Remove(addr)
		}
	}

	bytes, err := json.Marshal(stks)
	if err != nil {
		return nil, err
	}
	c.cache.Add(parent.Hash(), bytes)
	err = c.stakingDB.Commit(parent.Hash().Hex(), stks)
	if err != nil {
		return nil, err
	}

	return stks, nil
}

//[BERITH] 캐쉬나 db에서 stakingList를 불러오기 위한 메서드 생성
func (c *BSRR) getStakers(chain consensus.ChainReader, number uint64, hash common.Hash) (staking.Stakers, error) {
	var (
		list   staking.Stakers
		blocks []*types.Block
	)

	prevNum := number
	prevHash := hash

	//[BERITH] 입력받은 블록에서 가장가까운 StakingList를 찾는다.
	for list == nil {
		//[BERITH] cache에 저장된 StakingList를 찾은 경우
		if val, ok := c.cache.Get(prevHash); ok {
			bytes := val.([]byte)

			if err := json.Unmarshal(bytes, &list); err == nil {
				break
			}
			list = nil
			c.cache.Remove(prevHash)
		}

		//[BERITH] StakingList가 저장되지 않은 경우
		if prevNum == 0 {
			list = c.stakingDB.NewStakers()
			break
		}

		//[BERITH] DB에 저장된 StakingList를 찾은 경우

		var err error
		list, err = c.stakingDB.GetStakers(prevHash.Hex())
		if err == nil {
			break
		}
		list = nil

		block := chain.GetBlock(prevHash, prevNum)
		if block == nil {
			return nil, errors.New("unknown anccesstor")
		}

		blocks = append(blocks, block)
		prevNum--
		prevHash = block.ParentHash()
	}

	if len(blocks) == 0 {
		return list, nil
	}

	for i := 0; i < len(blocks)/2; i++ {
		blocks[i], blocks[len(blocks)-1-i] = blocks[len(blocks)-1-i], blocks[i]
	}

	err := c.checkBlocks(chain, list, blocks)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}
	c.cache.Add(hash, bytes)
	err = c.stakingDB.Commit(hash.Hex(), list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

//[BERITH] 블록을 확인하여 stakingList에 값을 세팅하기 위한 메서드 생성
func (c *BSRR) checkBlocks(chain consensus.ChainReader, stks staking.Stakers, blocks []*types.Block) error {
	if len(blocks) == 0 {
		return nil
	}

	for _, block := range blocks {
		if err := c.setStakersWithTxs(nil, chain, stks, block.Transactions(), block.Header()); err != nil {
			return err
		}
	}

	return nil
}

//[BERITH] 트랜잭션 배열을 조사하여 stakingList에 값을 세팅하기 위한 메서드 생성
func (c *BSRR) setStakersWithTxs(state *state.StateDB, chain consensus.ChainReader, stks staking.Stakers, txs []*types.Transaction, header *types.Header) error {
	number := header.Number

	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)

	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	prevState, err := chain.StateAt(parent.Root)

	if err != nil {
		return errMissingState
	}

	stkChanged := make(map[common.Address]bool)

	for _, tx := range txs {
		msg, err := tx.AsMessage(types.MakeSigner(chain.Config(), number))
		if err != nil {
			return err
		}

		//Main -> Main (일반 TX)
		if msg.Base() == types.Main && msg.Target() == types.Main {
			continue
		}

		//[BERITH] 2019-09-03
		//마지막 Staking의 블록번호가 저장되도록 수정
		//일반 Tx가 아닌 경우 Stake or Unstake
		if chain.Config().IsBIP1(number) && msg.Base() == types.Stake && msg.Target() == types.Main {
			stkChanged[msg.From()] = false
		} else if msg.Base() == types.Main && msg.Target() == types.Stake {
			stkChanged[msg.From()] = true
		}
	}

	for addr, isAdd := range stkChanged {
		if state != nil {
			point := big.NewInt(0)
			currentStkBal := state.GetStakeBalance(addr)
			if currentStkBal.Cmp(big.NewInt(0)) == 1 {
				currentStkBal = new(big.Int).Div(currentStkBal, common.UnitForBer)
				prevStkBal := new(big.Int).Div(prevState.GetStakeBalance(addr), common.UnitForBer)
				additionalStkBal := new(big.Int).Sub(currentStkBal, prevStkBal)
				currentBlock := header.Number
				lastStkBlock := new(big.Int).Set(state.GetStakeUpdated(addr))
				period := c.config.Period
				point = staking.CalcPointBigint(prevStkBal, additionalStkBal, currentBlock, lastStkBlock, period)
			}
			state.SetPoint(addr, point)
		}

		if isAdd {
			stks.Put(addr)
		} else {
			stks.Remove(addr)
		}

	}

	return nil
}

type signers []common.Address

func (s signers) signersMap() map[common.Address]struct{} {
	result := make(map[common.Address]struct{})
	for _, signer := range s {
		result[signer] = struct{}{}
	}
	return result
}

//[BERITH] 입력받은 블록넘버에, 블록생성이 가능한 계정의 목록을 반환하는 메서드.
// 1) [0, epoch number) -> genesis의 extra 데이터에서 추출 한 signers 반환
// 2) [epoch nunber ~ ) -> staking list 에서 추출 한 signers 반환
func (c *BSRR) getSigners(chain consensus.ChainReader, target *types.Header) (signers, error) {
	// extract signers from genesis block's extra data if block number equals to 0
	if target.Number.Cmp(big.NewInt(0)) == 0 {
		return c.getSignersFromExtraData(target)
	}

	// extract signers from genesis block's extra data if block number is less than epoch
	if target.Number.Cmp(big.NewInt(int64(c.config.Epoch))) < 0 {
		return c.getSignersFromExtraData(chain.GetHeaderByNumber(0))
	}

	// extract signers from staking list if block number is greater than or equals to epoch
	list, err := c.getStakers(chain, target.Number.Uint64(), target.Hash())
	if err != nil {
		return nil, errors.New("failed to get staking list")
	}

	result := list.AsList()
	if len(result) == 0 {
		return make([]common.Address, 0), nil
	}
	return result, nil
}

//[BERITH] getSignersFromExtraData extra data 필드로 부터 signers를 반환한다.
func (c *BSRR) getSignersFromExtraData(header *types.Header) (signers, error) {
	n := (len(header.Extra) - extraVanity - extraSeal) / common.AddressLength
	if n < 1 {
		return nil, errExtraSigners
	}

	signers := make([]common.Address, n)
	for i := 0; i < len(signers); i++ {
		copy(signers[i][:], header.Extra[extraVanity+i*common.AddressLength:])
	}
	return signers, nil
}

// [BERITH] getMaxMiningCandidates 주어진 스테이킹 리스트 수에서 블록을 생성 할 수 있는 후보자의 수를 반환한다.
func (c *BSRR) getMaxMiningCandidates(holders int) int {
	if holders == 0 {
		return 0
	}

	// (0,1) 범위는 모두 1
	t := int(math.Round(c.config.ForkFactor * float64(holders)))
	if t == 0 {
		t = 1
	}

	if t > selection.MAX_MINERS {
		t = selection.MAX_MINERS
	}
	return t
}

/*
[BERITH]
선출확율 반환 함수
*/
func (c *BSRR) getJoinRatio(stks staking.Stakers, address common.Address, hash common.Hash, blockNumber uint64, states *state.StateDB) (float64, error) {
	var total float64
	var n float64

	for _, stk := range stks.AsList() {
		point := float64(states.GetPoint(stk).Int64())
		if address == stk {
			n = point
		}
		total += point
	}

	if total == 0 {
		return 0, nil
	}

	return n / total, nil
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (c *BSRR) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "bsrr",
		Version:   "1.0",
		Service:   &API{chain: chain, bsrr: c},
		Public:    false,
	}}
}
