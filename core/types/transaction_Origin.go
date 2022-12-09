package types

import (
	"io"
	"math/big"
	"sync/atomic"

	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/crypto"
	"github.com/BerithFoundation/berith-chain/rlp"
)

type TransactionInterface interface {
	ChainId() *big.Int
	Protected() bool
	EncodeRLP(w io.Writer) error
	DecodeRLP(s *rlp.Stream) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(input []byte) error
	Data() []byte
	Gas() uint64
	GasPrice() *big.Int
	Value() *big.Int
	Nonce() uint64
	CheckNonce() bool
	Base() JobWallet
	Target() JobWallet
	To() *common.Address
	Hash() common.Hash
	Size() common.StorageSize
	AsMessage(s Signer) (Message, error)
	WithSignature(signer Signer, sig []byte) (*Transaction, error)
	Cost() *big.Int
	MainFee() *big.Int
	RawSignatureValues() (*big.Int, *big.Int, *big.Int)
}

type originTxdata struct {
	// From의 Nonce
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

func (o *originTxdata) MarshalJSON() ([]byte, error)     { return nil, nil }
func (o *originTxdata) UnmarshalJSON(input []byte) error { return nil }

type OriginTransaction struct {
	data originTxdata
	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

func (o *OriginTransaction) ChainId() *big.Int {
	return deriveChainId(o.data.V)
}

// Protected returns whether the transaction is protected from replay protection.
func (o *OriginTransaction) Protected() bool {
	return isProtectedV(o.data.V)
}

func (o *OriginTransaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &o.data)
}

// DecodeRLP implements rlp.Decoder
func (o *OriginTransaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&o.data)
	if err == nil {
		o.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

// MarshalJSON encodes the web3 RPC transaction format.
func (o *OriginTransaction) MarshalJSON() ([]byte, error) {
	hash := o.Hash()
	data := o.data
	data.Hash = &hash
	return data.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (o *OriginTransaction) UnmarshalJSON(input []byte) error {
	var dec originTxdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}

	withSignature := dec.V.Sign() != 0 || dec.R.Sign() != 0 || dec.S.Sign() != 0
	if withSignature {
		var V byte
		if isProtectedV(dec.V) {
			chainID := deriveChainId(dec.V).Uint64()
			V = byte(dec.V.Uint64() - 35 - 2*chainID)
		} else {
			V = byte(dec.V.Uint64() - 27)
		}
		if !crypto.ValidateSignatureValues(V, dec.R, dec.S, false) {
			return ErrInvalidSig
		}
	}

	*o = OriginTransaction{data: dec}
	return nil
}

func (o *OriginTransaction) Data() []byte       { return common.CopyBytes(o.data.Payload) }
func (o *OriginTransaction) Gas() uint64        { return o.data.GasLimit }
func (o *OriginTransaction) GasPrice() *big.Int { return new(big.Int).Set(o.data.Price) }
func (o *OriginTransaction) Value() *big.Int    { return new(big.Int).Set(o.data.Amount) }
func (o *OriginTransaction) Nonce() uint64      { return o.data.AccountNonce }
func (o *OriginTransaction) CheckNonce() bool   { return true }
func (o *OriginTransaction) Base() JobWallet    { return JobWallet(1) } //[Berith] Tx JobWallet Base
func (o *OriginTransaction) Target() JobWallet  { return JobWallet(1) } //[Berith] Tx JobWallet Target

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (o *OriginTransaction) To() *common.Address {
	if o.data.Recipient == nil {
		return nil
	}
	to := *o.data.Recipient
	return &to
}

// Hash hashes the RLP encoding of o.
// It uniquely identifies the transaction.
func (o *OriginTransaction) Hash() common.Hash {
	if hash := o.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(o)
	o.hash.Store(v)
	return v
}

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (o *OriginTransaction) Size() common.StorageSize {
	if size := o.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &o.data)
	o.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (o *OriginTransaction) AsMessage(s Signer) (Message, error) {
	msg := Message{
		nonce:      o.data.AccountNonce,
		gasLimit:   o.data.GasLimit,
		gasPrice:   new(big.Int).Set(o.data.Price),
		to:         o.data.Recipient,
		amount:     o.data.Amount,
		data:       o.data.Payload,
		checkNonce: true,
		base:       JobWallet(1),
		target:     JobWallet(1),
	}

	tx := &Transaction{
		data: txdata{
			AccountNonce: o.data.AccountNonce,
			Price:        o.data.Price,
			GasLimit:     o.data.GasLimit,
			Recipient:    o.data.Recipient,
			Amount:       o.data.Amount,
			Payload:      o.data.Payload,
			Base:         JobWallet(1),
			Target:       JobWallet(1),
			V:            o.data.V,
			R:            o.data.R,
			S:            o.data.S,
			Hash:         o.data.Hash},
		hash: o.hash,
		size: o.size,
		from: o.from,
	}
	var err error
	msg.from, err = Sender(s, tx)
	return msg, err
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (o *OriginTransaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	tx := &Transaction{
		data: txdata{
			AccountNonce: o.data.AccountNonce,
			Price:        o.data.Price,
			GasLimit:     o.data.GasLimit,
			Recipient:    o.data.Recipient,
			Amount:       o.data.Amount,
			Payload:      o.data.Payload,
			Base:         JobWallet(1),
			Target:       JobWallet(1),
			V:            o.data.V,
			R:            o.data.R,
			S:            o.data.S,
			Hash:         o.data.Hash},
		hash: o.hash,
		size: o.size,
		from: o.from,
	}
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}

// Cost returns amount + gasprice * gaslimit.
func (o *OriginTransaction) Cost() *big.Int {
	total := new(big.Int).Mul(o.data.Price, new(big.Int).SetUint64(o.data.GasLimit))
	total.Add(total, o.data.Amount)
	return total
}

func (o *OriginTransaction) MainFee() *big.Int {
	total := new(big.Int).Mul(o.data.Price, new(big.Int).SetUint64(o.data.GasLimit))
	return total
}

func (o *OriginTransaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return o.data.V, o.data.R, o.data.S
}

func NewOriginTransaction(tx *Transaction) *OriginTransaction {
	originTx := &OriginTransaction{
		data: originTxdata{
			AccountNonce: tx.data.AccountNonce,
			Price:        tx.data.Price,
			GasLimit:     tx.data.GasLimit,
			Recipient:    tx.data.Recipient,
			Amount:       tx.data.Amount,
			Payload:      tx.data.Payload,
			V:            tx.data.V,
			R:            tx.data.R,
			S:            tx.data.S,
			Hash:         tx.data.Hash},
		hash: tx.hash,
		size: tx.size,
		from: tx.from,
	}
	return originTx
}
