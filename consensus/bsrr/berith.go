package bsrr

import (
	"bytes"
	"errors"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"bitbucket.org/ibizsoftware/berith-chain/accounts"
	"bitbucket.org/ibizsoftware/berith-chain/berith/staking"
	"bitbucket.org/ibizsoftware/berith-chain/common"
	"bitbucket.org/ibizsoftware/berith-chain/consensus"
	"bitbucket.org/ibizsoftware/berith-chain/consensus/misc"
	"bitbucket.org/ibizsoftware/berith-chain/core/state"
	"bitbucket.org/ibizsoftware/berith-chain/core/types"
	"bitbucket.org/ibizsoftware/berith-chain/crypto"
	"bitbucket.org/ibizsoftware/berith-chain/crypto/sha3"
	"bitbucket.org/ibizsoftware/berith-chain/ethdb"
	"bitbucket.org/ibizsoftware/berith-chain/log"
	"bitbucket.org/ibizsoftware/berith-chain/params"
	"bitbucket.org/ibizsoftware/berith-chain/rlp"
	"bitbucket.org/ibizsoftware/berith-chain/rpc"
	"github.com/hashicorp/golang-lru"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory

	wiggleTime = 500 * time.Millisecond // Random delay (per signer) to allow concurrent signers
)

var (
	BlockReward = new(big.Int).Mul(big.NewInt(1e+18), big.NewInt(16))

	epochLength = uint64(30000) // Default number of blocks after which to checkpoint and reset the pending votes

	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.

	diffInTurn = big.NewInt(2) // Block difficulty for in-turn signatures
	diffNoTurn = big.NewInt(1) // Block difficulty for out-of-turn signatures
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

	// errWrongDifficulty is returned if the difficulty of a block doesn't match the
	// turn of the signer.
	errWrongDifficulty = errors.New("wrong difficulty")

	// ErrInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")

	// errUnauthorizedSigner is returned if a header is signed by a non-authorized entity.
	errUnauthorizedSigner = errors.New("unauthorized signer")
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

	rlp.Encode(hasher, []interface{}{
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

// ecrecover extracts the Ethereum account address from a signed header.
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

	// Recover the public key and the Ethereum address
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
	db     ethdb.Database     // Database to store and retrieve snapshot checkpoints

	//[BERITH] stakingDB clique 구조체에 추가
	stakingDB staking.DataBase //stakingList를 저장하는 DB
	cache     *lru.ARCCache    //stakingList를 저장하는 cache

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	proposals map[common.Address]bool // Current list of proposals we are pushing

	signer common.Address // Ethereum address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with
	lock   sync.RWMutex   // Protects the signer fields

	// The fields below are for testing only
	fakeDiff bool // Skip difficulty verifications
}

func New(config *params.BSRRConfig, db ethdb.Database) *BSRR {
	conf := config
	if conf.Epoch == 0 {
		conf.Epoch = epochLength
	}

	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)
	//[Berith] 캐쉬 인스턴스 생성및 사이즈 지정
	cache, _ := lru.NewARC(inmemorySnapshots)

	return &BSRR{
		config:     conf,
		db:         db,
		recents:    recents,
		signatures: signatures,
		cache:      cache,
		proposals:  make(map[common.Address]bool),
	}

}

func NewCliqueWithStakingDB(stakingDB staking.DataBase, config *params.BSRRConfig, db ethdb.Database) *BSRR {
	engine := New(config, db)
	engine.stakingDB = stakingDB
	// Synchronize the engine.config and chainConfig.
	return engine
}

// Author implements consensus.Engine, returning the Ethereum address recovered
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
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 {
		if header.Difficulty == nil || (header.Difficulty.Cmp(diffInTurn) != 0 && header.Difficulty.Cmp(diffNoTurn) != 0) {
			return errInvalidDifficulty
		}
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

// snapshot retrieves the authorization snapshot at a given point in time.
func (c *BSRR) snapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)
	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := c.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(c.config, c.signatures, c.db, hash); err == nil {
				log.Trace("Loaded voting snapshot from disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}
		// If we're at an checkpoint block, make a snapshot if it's known
		if number == 0 || (number%c.config.Epoch == 0 && chain.GetHeaderByNumber(number-1) == nil) {
			checkpoint := chain.GetHeaderByNumber(number)
			if checkpoint != nil {
				hash := checkpoint.Hash()

				signers := make([]common.Address, (len(checkpoint.Extra)-extraVanity-extraSeal)/common.AddressLength)
				for i := 0; i < len(signers); i++ {
					copy(signers[i][:], checkpoint.Extra[extraVanity+i*common.AddressLength:])
				}
				snap = newSnapshot(c.config, c.signatures, number, hash, signers)
				if err := snap.store(c.db); err != nil {
					return nil, err
				}
				log.Info("Stored checkpoint snapshot to disk", "number", number, "hash", hash)
				break
			}
		}
		// No snapshot for this header, gather the header and move backward
		var header *types.Header
		if len(parents) > 0 {
			// If we have explicit parents, pick from there (enforced)
			header = parents[len(parents)-1]
			if header.Hash() != hash || header.Number.Uint64() != number {
				return nil, consensus.ErrUnknownAncestor
			}
			parents = parents[:len(parents)-1]
		} else {
			// No explicit parents (or no more left), reach out to the database
			header = chain.GetHeader(hash, number)
			if header == nil {
				return nil, consensus.ErrUnknownAncestor
			}
		}
		headers = append(headers, header)
		number, hash = number-1, header.ParentHash
	}
	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}

	snap, err := snap.apply(chain, c.stakingDB, headers, c)
	if err != nil {
		return nil, err
	}
	c.recents.Add(snap.Hash, snap)

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if err = snap.store(c.db); err != nil {
			return nil, err
		}
		log.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}
	return snap, err
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
func (c *BSRR) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := c.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// Resolve the authorization key and check against signers
	signer, err := ecrecover(header, c.signatures)
	if err != nil {
		return err
	}
	if _, ok := snap.Signers[signer]; !ok {
		return errUnauthorizedSigner
	}

	// Ensure that the difficulty corresponds to the turn-ness of the signer
	if !c.fakeDiff {
		inturn := snap.inturn(header.Number.Uint64(), signer)
		if inturn && header.Difficulty.Cmp(diffInTurn) != 0 {
			return errWrongDifficulty
		}
		if !inturn && header.Difficulty.Cmp(diffNoTurn) != 0 {
			return errWrongDifficulty
		}
	}
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (c *BSRR) Prepare(chain consensus.ChainReader, header *types.Header) error {

	// If the block isn't a checkpoint, cast a random vote (good enough for now)
	//header.Coinbase = common.Address{}
	header.Nonce = types.BlockNonce{}

	number := header.Number.Uint64()
	// Assemble the voting snapshot to check which votes make sense
	snap, err := c.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	//[BERITH] stakingList를 확인할 블록번호를 논스로 지정하여 전파
	//블록넘버가 Epoch으로 나누어 떨어지지 않는경우 부모의 논스를 다시 전파
	header.Nonce = parent.Nonce

	// Set the correct difficulty
	header.Difficulty = CalcDifficulty(snap, c.signer)

	//[BERITH] 블록번호가 Epoch으로 나누어 떨어지는 경우 nonce값을 현재 블록의 번호로 변경한다.
	if number%c.config.Epoch == 0 {
		header.Nonce = types.EncodeNonce(number)
	}

	// Ensure the extra data has all it's components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]

	if number%c.config.Epoch == 0 {
		for _, signer := range snap.signers() {
			header.Extra = append(header.Extra, signer[:]...)
		}
	}
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
	// No block rewards in PoA, so the state remains as is and uncles are dropped
	accumulateRewards(chain.Config(), state, header, uncles)
	//[Berith] stakingList 처리 로직 추가
	stakingList, err := c.getStakingList(chain, header.Number.Uint64()-1, header.ParentHash)
	if err != nil {
		return nil, err
	}
	err = c.setStakingListWithTxs(chain, stakingList, txs, header.Number)
	if err != nil {
		return nil, err
	}
	stakingList.Finalize()

	//var snap *Snapshot
	_, err = c.snapshot(chain, header.Number.Uint64(), header.ParentHash, nil)
	if err != nil {
		return nil, err
	}

	//slashBadSigner(state, header, snap)


	stakingList.Print()
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
	// For 0-period chains, refuse to seal empty blocks (no reward but would spin sealing)
	if c.config.Period == 0 && len(block.Transactions()) == 0 {
		log.Info("Sealing paused, waiting for transactions")
		return nil
	}
	// Don't hold the signer fields for the entire sealing procedure
	c.lock.RLock()
	signer, signFn := c.signer, c.signFn
	c.lock.RUnlock()

	// Bail out if we're unauthorized to sign a block
	snap, err := c.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}
	if _, authorized := snap.Signers[signer]; !authorized {
		return errUnauthorizedSigner
	}

	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := time.Unix(header.Time.Int64(), 0).Sub(time.Now()) // nolint: gosimple
	if header.Difficulty.Cmp(diffNoTurn) == 0 {
		// It's not our turn explicitly to sign, delay it a bit
		wiggle := time.Duration(len(snap.Signers)/2+1) * wiggleTime
		delay += time.Duration(rand.Int63n(int64(wiggle)))

		log.Trace("Out-of-turn signing requested", "wiggle", common.PrettyDuration(wiggle))
	}
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
	snap, err := c.snapshot(chain, parent.Number.Uint64(), parent.Hash(), nil)
	if err != nil {
		return nil
	}
	return CalcDifficulty(snap, c.signer)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func CalcDifficulty(snap *Snapshot, signer common.Address) *big.Int {
	if snap.inturn(snap.Number+1, signer) {
		return new(big.Int).Set(diffInTurn)
	}
	return new(big.Int).Set(diffNoTurn)
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *BSRR) SealHash(header *types.Header) common.Hash {
	return sigHash(header)
}

// Close implements consensus.Engine. It's a noop for clique as there are no background threads.
func (c *BSRR) Close() error {
	return nil
}

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward and rewards for
// included uncles. The coinbase of each uncle block is also rewarded.
func accumulateRewards(config *params.ChainConfig, state *state.StateDB, header *types.Header, uncles []*types.Header) {
	number := header.Number.Uint64()
	if number <= 500 {
		return
	}

	//[BERITH]갯수 제한 코드 필요
	blockReward := BlockReward

	//state.AddStakeBalance(header.Coinbase, reward)
	state.AddRewardBalance(header.Coinbase, blockReward)
	//state.AddBalance(header.Coinbase, blockReward)
}

//[Berith] 제 차례에 블록을 쓰지 못한 마이너의 staking을 해제함
func slashBadSigner(state *state.StateDB, header *types.Header, snap *Snapshot) {
	number := header.Number.Uint64()
	signers := snap.signers()
	target := signers[number%uint64(len(signers))]

	if bytes.Compare(target.Bytes(), header.Coinbase.Bytes()) != 0 {
		state.AddBalance(target, state.GetStakeBalance(target))
		state.SetStaking(target, big.NewInt(0))
	}

}

//[Berith] 캐쉬나 db에서 stakingList를 불러오기 위한 메서드 생성
func (c *BSRR) getStakingList(chain consensus.ChainReader, number uint64, hash common.Hash) (staking.StakingList, error) {
	var (
		list   staking.StakingList
		blocks []*types.Block
	)

	for list == nil {
		if val, ok := c.cache.Get(hash); ok {
			bytes := val.([]byte)
			var err error
			list, err = staking.Decode(bytes)

			if err != nil {
				return nil, err
			}
			break
		}

		if number == 0 {
			list = c.stakingDB.NewStakingList()
			break
		}

		block := chain.GetBlock(hash, number)
		if block == nil {
			return nil, errors.New("unknown anccesstor")
		}

		blocks = append(blocks, block)
		number--
		hash = block.ParentHash()
	}

	for i := 0; i < len(blocks)/2; i++ {
		blocks[i], blocks[len(blocks)-1-i] = blocks[len(blocks)-1-i], blocks[i]
	}

	err := c.checkBlocks(chain, list, blocks)
	if err != nil {
		return nil, err
	}
	list.Finalize()
	list.Print()
	return list, nil

}

//[Berith] 블록을 확인하여 stakingList에 값을 세팅하기 위한 메서드 생성
func (c *BSRR) checkBlocks(chain consensus.ChainReader, stakingList staking.StakingList, blocks []*types.Block) error {
	if len(blocks) == 0 {
		return nil
	}

	for _, block := range blocks {
		c.setStakingListWithTxs(chain, stakingList, block.Transactions(), block.Number())
		number := block.NumberU64()
		snap, err := c.snapshot(chain, number, block.Hash(), nil)
		if err != nil {
			return err
		}
		signers := snap.signers()
		target := signers[number%uint64(len(signers))]
		coinbase := block.Header().Coinbase
		if !bytes.Equal(target.Bytes(), coinbase.Bytes()) {
			//stakingList.Delete(target)
		}
	}

	bytes, err := stakingList.Encode()
	if err != nil {
		return err
	}
	c.cache.Add(blocks[len(blocks)-1].Hash(), bytes)

	return nil
}

type stakingInfo struct {
	address     common.Address
	value       *big.Int
	blockNumber *big.Int
}

func (info stakingInfo) Address() common.Address { return info.address }
func (info stakingInfo) Value() *big.Int         { return info.value }
func (info stakingInfo) BlockNumber() *big.Int   { return info.blockNumber }

//[Berith] 트랜잭션 배열을 조사하여 stakingList에 값을 세팅하기 위한 메서드 생성
func (c *BSRR) setStakingListWithTxs(chain consensus.ChainReader, list staking.StakingList, txs []*types.Transaction, number *big.Int) error {
	for _, tx := range txs {
		msg, err := tx.AsMessage(types.MakeSigner(chain.Config(), number))
		if err != nil {
			return err
		}

		if msg.Base() != types.Stake && msg.Target() != types.Main {
			continue
		}

		//if !msg.Staking() && !bytes.Equal(msg.From().Bytes(), msg.To().Bytes()) {
		//	continue
		//}

		var info staking.StakingInfo
		info, err = list.GetInfo(msg.From())

		if err != nil {
			return err
		}

		blockNumber := number
		if info.BlockNumber().Cmp(blockNumber) > 0 {
			blockNumber = info.BlockNumber()
		}

		value := msg.Value()

		if msg.Base() == types.Stake && msg.Target() == types.Main {
			list.Delete(msg.From())
		} else if (msg.Base() == types.Main && msg.Target() ==types.Stake) && (msg.Base() == types.Reward && msg.Target() ==types.Stake) {
			input := stakingInfo{
				address:     msg.From(),
				value:       new(big.Int).Add(info.Value(), value),
				blockNumber: blockNumber,
			}

			list.SetInfo(input)
		}

		//if !msg.Staking() {
		//	// StopStaking()
		//	list.Delete(msg.From())
		//
		//} else if msg.Staking() && bytes.Equal(msg.From().Bytes(), msg.To().Bytes()){
		//	// Stake()
		//	input := stakingInfo{
		//		address:     msg.From(),
		//		value:       new(big.Int).Add(info.Value(), value),
		//		blockNumber: blockNumber,
		//	}
		//
		//	list.SetInfo(input)
		//}
	}
	return nil
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
