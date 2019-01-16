package bsrr

import (
	"bitbucket.org/ibizsoftware/berith-chain/accounts"
	"bitbucket.org/ibizsoftware/berith-chain/berith/stake"
	"bitbucket.org/ibizsoftware/berith-chain/common"
	"bitbucket.org/ibizsoftware/berith-chain/common/hexutil"
	"bitbucket.org/ibizsoftware/berith-chain/common/stakesort"

	"bitbucket.org/ibizsoftware/berith-chain/consensus"
	"bitbucket.org/ibizsoftware/berith-chain/core/rawdb"
	"bitbucket.org/ibizsoftware/berith-chain/core/types"
	"bitbucket.org/ibizsoftware/berith-chain/crypto"
	"bitbucket.org/ibizsoftware/berith-chain/crypto/sha3"
	"bitbucket.org/ibizsoftware/berith-chain/ethdb"
	"bitbucket.org/ibizsoftware/berith-chain/params"
	"bitbucket.org/ibizsoftware/berith-chain/rlp"
	"errors"
	"github.com/hashicorp/golang-lru"
	"math/big"
	"sync"
	"time"
)


type BSRR struct{
	config *params.BSRRConfig
	db ethdb.Database
	stakedb stake.DataBase
	stakelist *stakesort.Stakelist

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	proposals map[common.Address]bool // Current list of proposals we are pushing

	signer common.Address // Ethereum address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with
	lock   sync.RWMutex   // Protects the signer fields

	// The fields below are for testing only
	fakeDiff bool // Skip difficulty verifications



}
type SignerFn func(accounts.Account, []byte) ([]byte, error)
func New(config *params.BSRRConfig, db ethdb.Database,stakedb stake.DataBase) *BSRR{
	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)


	return &BSRR{
		config: config,
		db: db,
		stakedb: stakedb,
		recents:    recents,
		signatures: signatures,
		proposals:  make(map[common.Address]bool),
	}


}

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory

	wiggleTime = 500 * time.Millisecond // Random delay (per signer) to allow concurrent signers
)

var (
	epochLength = uint64(30000) // Default number of blocks after which to checkpoint and reset the pending votes

	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal

	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new signer
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a signer.

	uncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.

	diffInTurn = big.NewInt(2) // Block difficulty for in-turn signatures
	diffNoTurn = big.NewInt(1) // Block difficulty for out-of-turn signatures
)
var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errInvalidCheckpointBeneficiary is returned if a checkpoint/epoch transition
	// block has a beneficiary set to non-zeroes.
	errInvalidCheckpointBeneficiary = errors.New("beneficiary in checkpoint block non-zero")

	// errInvalidVote is returned if a nonce value is something else that the two
	// allowed constants of 0x00..0 or 0xff..f.
	errInvalidVote = errors.New("vote nonce not 0x00..0 or 0xff..f")

	// errInvalidCheckpointVote is returned if a checkpoint/epoch transition block
	// has a vote nonce set to non-zeroes.
	errInvalidCheckpointVote = errors.New("vote nonce in checkpoint block non-zero")

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

	// errMismatchingCheckpointSigners is returned if a checkpoint block contains a
	// list of signers different than the one the local node calculated.
	errMismatchingCheckpointSigners = errors.New("mismatching signer list on checkpoint block")

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

	// errRecentlySigned is returned if a header is signed by an authorized entity
	// that already signed a header recently, thus is temporarily not allowed to.
	errRecentlySigned = errors.New("recently signed")
)
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
func (b *BSRR)getStakingList(blocknumber *big.Int) (err error){
	hash :=rawdb.ReadCanonicalHash(b.db,blocknumber.Uint64())
	var stakelist stake.StakingList
	stakelist,err = stake.GetStakingMap(b.stakedb,blocknumber,hash)
	if err != nil{
		return err
	}
	b.stakelist = stakelist.GetRRList()

	return nil
}

func (b *BSRR) Author(header *types.Header) (common.Address, error){
	return ecrecover(header, b.signatures)

}
func (b *BSRR) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil{
		return errUnknownBlock
	}
	number := header.Number.Uint64()
	if header.Time.Cmp(big.NewInt(time.Now().Unix())) > 0{
		return consensus.ErrFutureBlock
	}




}
