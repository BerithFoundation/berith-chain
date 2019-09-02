/*
d8888b. d88888b d8888b. d888888b d888888b db   db
88  `8D 88'     88  `8D   `88'   `~~88~~' 88   88
88oooY' 88ooooo 88oobY'    88       88    88ooo88
88~~~b. 88~~~~~ 88`8b      88       88    88~~~88
88   8D 88.     88 `88.   .88.      88    88   88
Y8888P' Y88888P 88   YD Y888888P    YP    YP   YP

	  copyrights by ibizsoftware 2018 - 2019
*/

package bsrr

import (
	"bytes"
	"errors"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/BerithFoundation/berith-chain/rpc"

	"github.com/BerithFoundation/berith-chain/accounts"
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
	"github.com/gookit/color"

	"github.com/hashicorp/golang-lru"
)

const (
	inmemorySnapshots  = 128     // Number of recent vote snapshots to keep in memory
	inmemorySigners    = 128 * 3 // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096    // Number of recent block signatures to keep in memory

	//stakingInterval = 10
	wiggleTime = 500 * time.Millisecond // Random delay (per signer) to allow concurrent signers

)

var (
	RewardBlock  = big.NewInt(500)
	StakeMinimum = new(big.Int).Mul(big.NewInt(100000), big.NewInt(1e+18))
	SlashRound   = uint64(2)

	epochLength = uint64(360) // Default number of blocks after which to checkpoint and reset the pending votes

	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.

	//diffInTurn = big.NewInt(20000000) // Block difficulty for in-turn signatures
	//diffNoTurn = big.NewInt(10000000) // Block difficulty for out-of-turn signatures

	delays = []int{0, 1, 2, 3}
	groups = []int{3, 8, 15, 22}
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

	errStakingList = errors.New("not found staking list")
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

	proposals map[common.Address]bool // Current list of proposals we are pushing

	signer common.Address // Berith address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with
	lock   sync.RWMutex   // Protects the signer fields

	// The fields below are for testing only
	fakeDiff bool // Skip difficulty verifications
}

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

	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)
	//[Berith] 캐쉬 인스턴스 생성및 사이즈 지정
	cache, _ := lru.NewARC(inmemorySigners)

	return &BSRR{
		config:     conf,
		db:         db,
		recents:    recents,
		signatures: signatures,
		cache:      cache,
		proposals:  make(map[common.Address]bool),
	}

}

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
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	//if number > 0 {
	//	if header.Difficulty == nil || header.Difficulty.Uint64() > diffInTurn.Uint64() {
	//		return errInvalidDifficulty
	//	}
	//	if header.Difficulty == nil || (header.Difficulty.Cmp(diffInTurn) != 0 && header.Difficulty.Cmp(diffNoTurn) != 0) {
	//		return errInvalidDifficulty
	//	}
	//}
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
func (c *BSRR) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number
	if number.Uint64() == 0 {
		return errUnknownBlock
	}
	//signers := c.getSigners(chain, header)

	// Resolve the authorization key and check against signers
	// signer, err := ecrecover(header, c.signatures)
	// if err != nil {
	// 	return err
	// }
	// if _, ok := signers.signersMap()[signer]; !ok {
	// 	return errUnauthorizedSigner
	// }

	// if !c.fakeDiff {
	// 	inturn := signers[(header.Number.Uint64()%c.config.Epoch)%uint64(len(signers))] == signer
	// 	if inturn && header.Difficulty.Cmp(diffInTurn) != 0 {
	// 		return errWrongDifficulty
	// 	}
	// 	if !inturn && header.Difficulty.Cmp(diffNoTurn) != 0 {
	// 		return errWrongDifficulty
	// 	}
	// }
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (c *BSRR) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// If the block isn't a checkpoint, cast a random vote (good enough for now)
	//header.Coinbase = common.Address{}
	header.Nonce = types.BlockNonce{}

	number := header.Number.Uint64()

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	//[BERITH] stakingList를 확인할 블록번호를 논스로 지정하여 전파
	//블록넘버가 Epoch으로 나누어 떨어지지 않는경우 부모의 논스를 다시 전파
	header.Nonce = parent.Nonce

	// Set the correct difficulty
	diff, rank := c.calcDifficultyAndRank(c.signer, chain, number-1, parent)

	if rank > staking.MAX_MINERS {
		return errUnauthorizedSigner
	}

	header.Difficulty = diff

	//[BERITH] 블록번호가 Epoch으로 나누어 떨어지는 경우 nonce값을 현재 블록의 번호로 변경한다.
	if number%c.config.Epoch == 0 {
		header.Nonce = types.EncodeNonce(number)
	}

	// Ensure the extra data has all it's components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]

	header.Extra = append(header.Extra, make([]byte, extraSeal)...)

	// Mix digest is reserved for now, set to empty
	header.MixDigest = common.Hash{}

	// Ensure the timestamp has the correct delay
	//parent := chain.GetHeader(header.ParentHash, number-1)
	// if parent == nil {
	// 	return consensus.ErrUnknownAncestor
	// }
	header.Time = new(big.Int).Add(parent.Time, new(big.Int).SetUint64(c.config.Period))
	if header.Time.Int64() < time.Now().Unix() {
		header.Time = big.NewInt(time.Now().Unix())
	}
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (c *BSRR) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	//[Berith] stakingList 처리 로직 추가
	stakingList, err := c.getStakingList(chain, header.Number.Uint64()-1, header.ParentHash)
	if err != nil {
		return nil, errStakingList
	}

	font := color.Yellow
	if bytes.Compare(header.Coinbase.Bytes(), c.signer.Bytes()) == 0 {
		font = color.Green
	}
	//stakingList.Print()
	font.Println("##############[FINALIZE]##############")
	font.Println("NUMBER : ", header.Number.String())
	font.Println("HASH : ", header.Hash().Hex())
	font.Println("COINBASE : ", header.Coinbase.Hex())
	font.Println("DIFFICULTY : ", header.Difficulty.String())
	font.Println("UNCLES : ", header.UncleHash.Hex())
	font.Println("######################################")

	if header.Coinbase != common.HexToAddress("0") {
		var signers signers
		//Diff
		epoch := chain.Config().Bsrr.Epoch
		targetNumber := header.Number.Uint64() - epoch

		signers, err := c.getSigners(chain, header.Number.Uint64()-1, targetNumber, header.ParentHash)
		if err != nil {
			return nil, errUnauthorizedSigner
		}

		signerMap := signers.signersMap()
		if _, ok := signerMap[header.Coinbase]; !ok {
			return nil, errUnauthorizedSigner
		}

		parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
		predicted, rank := c.calcDifficultyAndRank(header.Coinbase, chain, 0, parent)
		if rank > staking.MAX_MINERS {
			return nil, errUnauthorizedSigner
		}
		font = color.Blue
		font.Println("Remote :: " + header.Difficulty.String() + "\tLocal :: " + predicted.String())
		if predicted.Cmp(header.Difficulty) != 0 {
			return nil, errInvalidDifficulty
		}
	}

	err = c.setStakingListWithTxs(state, chain, stakingList, txs, header)
	if err != nil {
		return nil, errStakingList
	}

	//Reward 보상
	c.accumulateRewards(chain, state, header)

	//상태값 적용
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
	epoch := chain.Config().Bsrr.Epoch
	targetNumber := header.Number.Uint64() - epoch
	signers, err := c.getSigners(chain, header.Number.Uint64()-1, targetNumber, header.ParentHash)

	//signers, err := c.getSigners(chain, header.Number.Uint64()-1, header.ParentHash)
	if err != nil {
		return err
	}

	if _, authorized := signers.signersMap()[signer]; !authorized {
		return errUnauthorizedSigner
	}

	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := time.Unix(header.Time.Int64(), 0).Sub(time.Now()) // nolint: gosimple
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	_, rank := c.calcDifficultyAndRank(header.Coinbase, chain, 0, parent)

	additionalDelay := -1

	for i := 0; i < len(groups); i++ {
		if rank <= groups[i] {
			additionalDelay = delays[i]
			break
		}
	}

	if additionalDelay == -1 {
		return errUnauthorizedSigner
	}

	delay += time.Duration(additionalDelay) * time.Second

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
	diff, _ := c.calcDifficultyAndRank(c.signer, chain, time, parent)
	return diff
}
func (c *BSRR) calcDifficultyAndRank(signer common.Address, chain consensus.ChainReader, time uint64, parent *types.Header) (*big.Int, int) {

	target := parent
	targetNumber := new(big.Int).Sub(parent.Number, big.NewInt(int64(c.config.Epoch)))
	for target.Number.Cmp(big.NewInt(0)) > 0 && target.Number.Cmp(targetNumber) > 0 {
		target = chain.GetHeader(target.ParentHash, target.Number.Uint64()-1)
	}
	if target.Number.Cmp(big.NewInt(0)) <= 0 {
		return big.NewInt(1234), 1
	}

	list, err := c.getStakingList(chain, target.Number.Uint64(), target.Hash())

	if err != nil {
		return big.NewInt(0), staking.MAX_MINERS + 1
	}

	diff, rank, reordered := list.GetDifficultyAndRank(signer, target.Number.Uint64(), c.config.Period)
	if reordered {
		bytes, _ := list.Encode()
		c.cache.Add(target.Hash(), bytes)

		if target.Number.Uint64()%c.config.Period == 0 {
			c.stakingDB.Commit(target.Hash().Hex(), list)
		}
	}

	return diff, rank
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *BSRR) SealHash(header *types.Header) common.Hash {
	return sigHash(header)
}

// Close implements consensus.Engine. It's a noop for clique as there are no background threads.
func (c *BSRR) Close() error {
	return nil
}

func getReward(config *params.ChainConfig, header *types.Header) *big.Int {
	number := header.Number.Uint64()
	if number < config.Bsrr.Rewards.Uint64() {
		return big.NewInt(0)
	}

	d := float64(config.Bsrr.Period) / 10
	n := float64(number) * d

	var z float64 = 0
	if n <= 3150000 {
		z = 5
	}

	re := (26 - math.Round(n/(7370000))*0.5 + z) * d
	if re <= 0 {
		re = 0

		return big.NewInt(0)
	} else {
		temp := re * 1e+10
		return new(big.Int).Mul(big.NewInt(int64(temp)), big.NewInt(1e+8))
	}
}

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward and rewards for
// included uncles. The coinbase of each uncle block is also rewarded.
func (c *BSRR) accumulateRewards(chain consensus.ChainReader, state *state.StateDB, header *types.Header) {

	config := chain.Config()
	state.AddBehindBalance(header.Coinbase, header.Number, getReward(config, header))

	//과거 시점의 블록 생성자 가져온다.
	targetNumber := header.Number.Uint64() - config.Bsrr.SlashRound
	signers, err := c.getSigners(chain, header.Number.Uint64()-1, targetNumber, header.ParentHash)
	if err != nil {
		return
	}

	//all node block result
	for _, addr := range signers {
		behind, err := state.GetFirstBehindBalance(addr)
		if err != nil {
			continue
		}

		target := new(big.Int).Add(behind.Number, new(big.Int).SetUint64(config.Bsrr.SlashRound))
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

//[Berith] 제 차례에 블록을 쓰지 못한 마이너의 staking을 해제함
func (c *BSRR) slashBadSigner(chain consensus.ChainReader, header *types.Header, list staking.StakingList, state *state.StateDB) error {

	epoch := chain.Config().Bsrr.Epoch
	targetNumber := header.Number.Uint64() - epoch
	signers, err := c.getSigners(chain, header.Number.Uint64(), targetNumber, header.Hash())
	//signers, err := c.getSigners(chain, header.Number.Uint64()-1, header.ParentHash)
	if err != nil {
		return err
	}

	signerMap := make(map[common.Address]bool)

	for _, val := range signers {
		signerMap[val] = true
	}

	miners := list.GetMiners()

	for k, _ := range signerMap {
		_, ok := miners[k]
		if !ok {
			if state != nil {
				state.AddBalance(k, state.GetStakeBalance(k))
				state.SetStaking(k, big.NewInt(0))
			}
			info, err := list.GetInfo(k)
			if err != nil {
				return err
			}
			list.SetInfo(&stakingInfo{
				address:     info.Address(),
				value:       big.NewInt(0),
				blockNumber: info.BlockNumber(),
				reward:      info.Reward(),
			})
		}
	}

	list.InitMiner()

	return nil

}

//[Berith] 캐쉬나 db에서 stakingList를 불러오기 위한 메서드 생성
func (c *BSRR) getStakingList(chain consensus.ChainReader, number uint64, hash common.Hash) (staking.StakingList, error) {
	var (
		list   staking.StakingList
		blocks []*types.Block
	)

	prevNum := number
	prevHash := hash

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

		if prevNum == 0 {
			list = c.stakingDB.NewStakingList()
			list.SetTarget(prevHash)
			break
		}

		if prevNum%c.config.Epoch == 0 {
			var err error
			list, err = c.stakingDB.GetStakingList(prevHash.Hex())

			if err == nil {
				break
			}

			list = nil
		}

		block := chain.GetBlock(prevHash, prevNum)
		if block == nil {
			return nil, errors.New("unknown anccesstor")
		}

		blocks = append(blocks, block)
		prevNum--
		prevHash = block.ParentHash()
	}

	for i := 0; i < len(blocks)/2; i++ {
		blocks[i], blocks[len(blocks)-1-i] = blocks[len(blocks)-1-i], blocks[i]
	}

	if len(blocks) > 0 {
		list.ClearTable()
	}

	list = list.Copy()

	err := c.checkBlocks(chain, list, blocks)
	if err != nil {
		return nil, err
	}

	list.Sort()

	bytes, err := list.Encode()
	if err != nil {
		return nil, err
	}
	c.cache.Add(number, bytes)

	if number%c.config.Period == 0 {
		c.stakingDB.Commit(hash.Hex(), list)
	}

	return list, nil

}

//[Berith] 블록을 확인하여 stakingList에 값을 세팅하기 위한 메서드 생성
func (c *BSRR) checkBlocks(chain consensus.ChainReader, stakingList staking.StakingList, blocks []*types.Block) error {
	if len(blocks) == 0 {
		return nil
	}

	for _, block := range blocks {
		c.setStakingListWithTxs(nil, chain, stakingList, block.Transactions(), block.Header())
		if block.NumberU64()%c.config.Epoch == 0 {
			stakingList.SetTarget(block.Hash())
		}
	}

	return nil
}

type stakingInfo struct {
	address     common.Address
	value       *big.Int
	blockNumber *big.Int
	reward      *big.Int
}

func (info stakingInfo) Address() common.Address { return info.address }
func (info stakingInfo) Value() *big.Int         { return info.value }
func (info stakingInfo) BlockNumber() *big.Int   { return info.blockNumber }
func (info stakingInfo) Reward() *big.Int        { return info.reward }

//[Berith] 트랜잭션 배열을 조사하여 stakingList에 값을 세팅하기 위한 메서드 생성
func (c *BSRR) setStakingListWithTxs(state *state.StateDB, chain consensus.ChainReader, list staking.StakingList, txs []*types.Transaction, header *types.Header) error {
	number := header.Number
	for _, tx := range txs {
		msg, err := tx.AsMessage(types.MakeSigner(chain.Config(), number))
		if err != nil {
			return err
		}

		//Main -> Main (일반 TX)
		if msg.Base() == types.Main && msg.Target() == types.Main {
			continue
		}

		//Reward -> Main
		// if (msg.Base() == types.Reward && msg.Target() == types.Main) &&
		// 	bytes.Equal(msg.From().Bytes(), msg.To().Bytes()) {
		// 	continue
		// }

		var info staking.StakingInfo
		info, err = list.GetInfo(msg.From())

		if err != nil {
			return err
		}

		value := new(big.Int).Set(info.Value())
		reward := new(big.Int).Set(info.Reward())

		//Stake
		if msg.Target() == types.Stake {
			value.Add(value, msg.Value())
			//add point
		}

		//Unstake
		if msg.Base() == types.Stake && msg.Target() == types.Main {
			value.Set(big.NewInt(0))
			//reset point
		}

		blockNumber := number
		if info.BlockNumber().Cmp(blockNumber) > 0 {
			blockNumber = info.BlockNumber()
		}

		input := stakingInfo{
			address:     msg.From(),
			value:       value,
			blockNumber: blockNumber,
			reward:      reward,
		}

		list.SetInfo(input)
	}

	info, err := list.GetInfo(header.Coinbase)

	if err != nil {
		return err
	}

	if info.Value().Cmp(big.NewInt(0)) > 0 {

		input := stakingInfo{
			address:     info.Address(),
			value:       info.Value(),
			blockNumber: info.BlockNumber(),
			reward:      new(big.Int).Add(info.Reward(), getReward(chain.Config(), header)),
		}

		list.SetInfo(input)
	}

	// list.SetMiner(header.Coinbase)
	// sr := c.config.SlashRound
	// if header.Number.Uint64()%(sr*c.config.Epoch) == 0 {
	// 	return c.slashBadSigner(chain, header, list, state)
	// }
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

func (c *BSRR) getSigners(chain consensus.ChainReader, number, targetNumber uint64, hash common.Hash) (signers, error) {
	checkpoint := chain.GetHeaderByNumber(0)
	signers := make([]common.Address, (len(checkpoint.Extra)-extraVanity-extraSeal)/common.AddressLength)
	for i := 0; i < len(signers); i++ {
		copy(signers[i][:], checkpoint.Extra[extraVanity+i*common.AddressLength:])
	}

	target := chain.GetHeader(hash, number)
	//targetNumber := number - c.config.Epoch

	if targetNumber <= 0 {
		return signers, nil
	}
	for target.Number.Uint64() > 0 && target.Number.Uint64() > targetNumber {
		target = chain.GetHeader(target.ParentHash, target.Number.Uint64()-1)
		if target == nil {
			return nil, errors.New("invalid ancestor")
		}
	}

	list, err := c.getStakingList(chain, target.Number.Uint64(), target.Hash())

	if err != nil {
		return nil, errors.New("Failed to get staking list")
	}

	list.Sort()

	result := list.ToArray()

	if len(result) <= 0 {
		return signers, nil
	}

	return result, nil

}

func (c *BSRR) getJoinRatio(stakingList *staking.StakingList, address common.Address, blockNumber uint64) (float64, error) {
	roi := (*stakingList).GetJoinRatio(address, blockNumber, c.config.Period)

	return roi, nil
}

func reward(number float64) float64 {
	up := 5.5 * 100 * math.Pow(10, 7.2)
	down := number + math.Pow(10, 7.6)

	y := up/down - 60.0

	if y < 0 {
		return float64(0)
	}
	return y
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
